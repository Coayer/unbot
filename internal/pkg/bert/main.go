package bert

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"log"
	"math"
)

var Model = loadModel()

const DIM int = 384

//loadModel loads TensorFlow SavedModel from disk
func loadModel() *tf.SavedModel {
	model, err := tf.LoadSavedModel("data/internal/qa_model", []string{"serve"}, nil)

	if err != nil {
		log.Fatal(err)
	}

	return model
}

//askBert is used to ask the BERT model a question
func AskBert(query string, context string) string {
	queryTokens, contextTokens := tokenize(query), tokenize(context)

	//context is reduced to part of original
	startToken, endToken, contextTokens := longContextForward(queryTokens, contextTokens)

	log.Println(contextTokens)

	/*start and end tokens are the start and end indices of the answer span of the context
	the resulting indices are relative to contextTokens, and not the BERT output vectors
	this is because contextTokens contains the encoded tokens, whereas BERT outputs are only weights*/
	startToken = startToken - len(queryTokens) - 2

	if startToken < 0 {
		return "Can't find answer"
	}

	endToken = endToken - len(queryTokens) - 1

	if endToken < 0 {
		return "Can't find answer"
	}

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
func longContextForward(queryTokens []string, contextTokens []string) (int, int, []string) {
	chunkSize := DIM - len(queryTokens) - 3 //control tokens, chunks needed to reduce context dimension
	var startToken, endToken int
	var bestContextTokens []string
	highestConfidence := float32(math.Inf(-1)) //highest sum of BERT's confidence of answer span

	for i := 0; i < len(contextTokens); i += chunkSize {
		var tempContextTokens []string

		//used for last chunk where there aren't enough tokens for full chunk
		if len(contextTokens)-i < chunkSize {
			tempContextTokens = contextTokens[i:]
		} else {
			tempContextTokens = contextTokens[i : i+chunkSize]
		}
		//start and end vectors for current subcontext
		tempStartVector, tempEndVector := forward(makeInputIDs(queryTokens, tempContextTokens), makeAttentionMask(len(queryTokens), len(tempContextTokens)))

		tempStartToken := argmax(tempStartVector)
		tempEndToken := argmax(tempEndVector[tempStartToken:]) + tempStartToken
		tempContextScore := tempStartVector[tempStartToken] + tempEndVector[tempEndToken] //confidence score indicating how much BERT thinks there is an answer

		log.Printf("Confidence: %f", tempContextScore)

		if tempContextScore > highestConfidence {
			startToken, endToken = tempStartToken, tempEndToken
			bestContextTokens = tempContextTokens
			highestConfidence = tempContextScore
		}
	}

	return startToken, endToken, bestContextTokens
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
