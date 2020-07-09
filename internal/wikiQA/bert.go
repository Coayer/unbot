package wikiQA

import (
	"log"
	"math"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

var Model = loadModel()

const dim int = 384

//loadModel loads TensorFlow SavedModel from disk
func loadModel() *tf.SavedModel {
	model, err := tf.LoadSavedModel("internal/wikiQA/model", []string{"serve"}, nil)

	if err != nil {
		log.Fatal(err)
	}

	return model
}

//askBert is used to ask the BERT model a question
func askBert(query string, context string) string {
	queryTokens := tokenize(query)
	contextTokens, piecesPerToken := tokenizePiecesCount(context)

	startVector, endVector := make([]float32, dim), make([]float32, dim) //BERT output vectors

	//BERT has limited input vector dimension, so need to know which part of the full context was used so it can be decoded later
	subContextStart, subContextEnd := 0, len(contextTokens)

	if len(queryTokens)+len(contextTokens)+3 > dim { //+3 for control tokens in BERT input
		log.Println("Long context")
		startVector, endVector, subContextStart, subContextEnd = longContextForward(queryTokens, contextTokens)
	} else {
		log.Println("Short context")
		startVector, endVector = forward(makeInputIDs(queryTokens, contextTokens), makeAttentionMask(len(queryTokens), len(contextTokens)))
	}

	log.Println(subContextStart, subContextEnd)
	log.Println(contextTokens[subContextStart:subContextEnd])

	//makes probability distributions for start and end positions of answer
	startVector = softmax(startVector)
	endVector = softmax(endVector)

	/*start and end tokens are the start and end indicies of the answer span of the context
	the resulting indicies are relative to contextTokens, and not the BERT output vectors
	this is because contextTokens contains the encoded tokens, whereas BERT outputs are only weights*/
	startToken := argmax(startVector) - len(queryTokens) - 1
	endToken := argmax(endVector[startToken:]) + startToken - len(queryTokens)

	/*at this point, start and end tokens are for the wordpieces, which would give answer of "129 ##5" instead of "1295"
	the wordpiece indicies are converted to "regular" tokenization (whole word tokenization) indicies
	this prevents words from being cut off and decodes the wordpiece tokens to their full words*/
	startToken = matchingTokenPosition(piecesPerToken[subContextStart:subContextEnd], startToken)
	endToken = matchingTokenPosition(piecesPerToken[subContextStart:subContextEnd], endToken)

	log.Printf("Found answer span: %d --> %d", startToken, endToken)

	answer := strings.Join(baseTokenize(context)[subContextStart:subContextEnd][startToken:endToken], " ")
	log.Println("Query: " + query)
	log.Println("Answer: " + answer)
	return answer
}

//forward runs the model
func forward(inputIDs []int32, attentionMask []int32) ([]float32, []float32) {
	inputIDsTensor, _ := tf.NewTensor(inputIDs)
	attentionMaskTensor, _ := tf.NewTensor(attentionMask)

	results, err := Model.Session.Run(
		map[tf.Output]*tf.Tensor{
			Model.Graph.Operation("serving_default_attention_mask").Output(0): attentionMaskTensor,
			Model.Graph.Operation("serving_default_input_ids").Output(0):      inputIDsTensor,
		},
		[]tf.Output{
			Model.Graph.Operation("StatefulPartitionedCall").Output(0),
			Model.Graph.Operation("StatefulPartitionedCall").Output(1),
		},
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	return results[0].Value().([][]float32)[0], results[1].Value().([][]float32)[0]
}

//longContextForward handles forwarding for contexts longer than dimension of BERT
func longContextForward(queryTokens []string, contextTokens []string) ([]float32, []float32, int, int) {
	chunkSize := dim - len(queryTokens) - 3 //control tokens

	highestStartEnd := float32(-1000) //highest sum of BERT's confidence of answer span
	startVector, endVector := make([]float32, dim), make([]float32, dim)
	subContextStart, subContextEnd := 0, len(contextTokens) //context from which the highest confidence of answer was found

	for i := 0; i < len(contextTokens); i += chunkSize {
		var tempContext []string

		//reduces size of context to fit in BERT input
		if len(contextTokens)-i < chunkSize {
			//used for last chunk where there aren't enough tokens for full chunk
			tempContext = contextTokens[i : len(contextTokens)-1]
		} else {
			tempContext = contextTokens[i : i+chunkSize]
		}

		//start and end vectors for current subcontext
		tempStart, tempEnd := forward(makeInputIDs(queryTokens, tempContext), makeAttentionMask(len(queryTokens), len(tempContext)))

		startToken := argmax(tempStart)
		endToken := argmax(tempEnd[startToken:]) + startToken
		tempContextScore := tempStart[startToken] + tempEnd[endToken] //confidence score indicating how much BERT thinks there is an answer

		log.Printf("Confidence: %f", tempContextScore)

		if tempContextScore > highestStartEnd {
			startVector, endVector = tempStart, tempEnd

			if len(contextTokens)-i < chunkSize {
				subContextStart, subContextEnd = i, len(contextTokens)-1
			} else {
				subContextStart, subContextEnd = i, i+chunkSize
			}

			highestStartEnd = tempContextScore
		}
	}

	return startVector, endVector, subContextStart, subContextEnd
}

//makeInputIDs encodes query and context into BERT readable format
func makeInputIDs(query []string, context []string) []int32 {
	result := make([]int32, 0, dim)
	result = append(result, 102) //[CLS]
	result = append(result, encodeTokens(query)...)
	result = append(result, 103) //[SEP]
	result = append(result, encodeTokens(context)...)
	result = append(result, 103) //[SEP]

	pad := make([]int32, dim-len(result))
	for i := range pad {
		pad[i] = 1 //[PAD]
	}

	return append(result, pad...)
}

//encodeTokens converts a list of tokens into BERT encodings
func encodeTokens(tokens []string) []int32 {
	result := make([]int32, len(tokens))

	for i, token := range tokens {
		if vocab.vocabMap[token] {
			result[i] = vocab.encodeMap[token]
		} else {
			result[i] = 101 //[UNK] in vocab
		}
	}

	return result
}

//makeAttentionMask creates attention mask given query and context lengths
func makeAttentionMask(queryLength int, contextLength int) []int32 {
	mask := make([]int32, dim)

	for i := 0; i < queryLength+contextLength+3; i++ {
		mask[i] = 1
	}

	return mask
}

//softmax performs softmax function on vector
func softmax(vector []float32) []float32 {
	for i, weight := range vector {
		vector[i] = float32(math.Pow(2.7182818284, float64(weight)))
	}

	total := sum32(vector)

	for i, weight := range vector {
		vector[i] = weight / total
	}

	return vector
}

//argmax finds index with the highest value
func argmax(vector []float32) int {
	maxValue := float32(0)
	maxIndex := 0

	for i, value := range vector {
		if value > maxValue {
			maxIndex = i
			maxValue = value
		}
	}

	return maxIndex
}

//sum32 calculates a sum of a set of float32 values
func sum32(values []float32) float32 {
	result := float32(0)

	for _, value := range values {
		result += value
	}

	return result
}

//matchingTokenPosition finds the index of a basic tokenized token given an index of a wordpiece token (see askBert)
func matchingTokenPosition(piecesPerToken []int, piecePosition int) int {
	total := 0

	for i, pieceCount := range piecesPerToken {
		total += pieceCount
		if total >= piecePosition {
			return i
		}
	}

	return -1 //shouldn't get to this
}
