package wikiQA

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"log"
	"math"
)

var Model = loadModel()

const DIM int = 384

//loadModel loads TensorFlow SavedModel from disk
func loadModel() *tf.SavedModel {
	model, err := tf.LoadSavedModel("data/qa_model", []string{"serve"}, nil)

	if err != nil {
		log.Fatal(err)
	}

	return model
}

//askBert is used to ask the BERT model a question
func askBert(query string, context string) string {
	queryTokens, contextTokens := tokenize(query), tokenize(context)
	startVector, endVector := make([]float32, DIM), make([]float32, DIM) //BERT output vectors

	if len(queryTokens)+len(contextTokens)+3 > DIM { //+3 for control tokens in BERT input, limited BERT dimension so need check
		log.Println("Long context")
		//context is reduced to part of original
		startVector, endVector, contextTokens = longContextForward(queryTokens, contextTokens)
	} else {
		log.Println("Short context")
		startVector, endVector = forward(makeInputIDs(queryTokens, contextTokens), makeAttentionMask(len(queryTokens), len(contextTokens)))
	}

	log.Println(contextTokens)

	//makes probability distributions for start and end positions of answer
	startVector = softmax(startVector)
	endVector = softmax(endVector)

	/*start and end tokens are the start and end indices of the answer span of the context
	the resulting indices are relative to contextTokens, and not the BERT output vectors
	this is because contextTokens contains the encoded tokens, whereas BERT outputs are only weights*/
	startToken := argmax(startVector) - len(queryTokens) - 2

	if startToken < 0 {
		return "Can't find answer"
	}

	endToken := argmax(endVector[startToken:]) + startToken - len(queryTokens) - 1

	log.Printf("Found answer span: %d --> %d", startToken, endToken)

	if endToken-startToken > 20 {
		endToken = startToken + 20
		log.Println("Clipping answer")
	}

	return tokensToEnglish(contextTokens[startToken:endToken])
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
func longContextForward(queryTokens []string, contextTokens []string) ([]float32, []float32, []string) {
	chunkSize := DIM - len(queryTokens) - 3 //control tokens, chunks needed to reduce context dimension
	startVector, endVector := make([]float32, DIM), make([]float32, DIM)
	highestConfidence := float32(math.Inf(-1)) //highest sum of BERT's confidence of answer span

	for i := 0; i < len(contextTokens); i += chunkSize {
		var tempContextTokens []string

		//used for last chunk where there aren't enough tokens for full chunk
		if len(contextTokens)-i < chunkSize {
			tempContextTokens = contextTokens[i : len(contextTokens)-1]
		} else {
			tempContextTokens = contextTokens[i : i+chunkSize]
		}

		//start and end vectors for current subcontext
		tempStart, tempEnd := forward(makeInputIDs(queryTokens, tempContextTokens), makeAttentionMask(len(queryTokens), len(tempContextTokens)))

		startToken := argmax(tempStart)
		endToken := argmax(tempEnd[startToken:]) + startToken
		tempContextScore := tempStart[startToken] + tempEnd[endToken] //confidence score indicating how much BERT thinks there is an answer

		log.Printf("Confidence: %f", tempContextScore)

		if tempContextScore > highestConfidence {
			startVector, endVector = tempStart, tempEnd
			contextTokens = tempContextTokens
			highestConfidence = tempContextScore
		}
	}

	return startVector, endVector, contextTokens
}

//makeInputIDs encodes query and context into BERT readable format
func makeInputIDs(query []string, context []string) []int32 {
	result := make([]int32, 0, DIM)
	result = append(result, 102) //[CLS]
	result = append(result, encodeTokens(query)...)
	result = append(result, 103) //[SEP]
	result = append(result, encodeTokens(context)...)
	result = append(result, 103) //[SEP]

	pad := make([]int32, DIM-len(result))
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
	mask := make([]int32, DIM)

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
	maxValue := float32(math.Inf(-1))
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
