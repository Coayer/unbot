package memory

import (
	"bufio"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

const DIM = 50

var glove = readGloVe()

func sentence2Vec(sentence string) []float32 {
	sentence = utils.RemoveStopWords(strings.ToLower(sentence))
	tokens := utils.BaseTokenize(sentence)

	tokenVectors := make([][]float32, len(tokens))
	for i := 0; i < len(tokens); i++ {
		tokenVectors[i] = (*glove)[tokens[i]]
	}

	vector := make([]float32, DIM)
	for j := 0; j < DIM; j++ {
		var val float32
		for i := 0; i < len(tokenVectors); i++ {
			val += tokenVectors[i][j]
		}
		vector[j] = val
	}
	return vector
}

func readGloVe() *map[string][]float32 {
	log.Println("Loading GloVe embeddings")

	glove := make(map[string][]float32)
	file, err := os.Open("data/glove50d.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")

		embedding := make([]float32, DIM)
		for i := 0; i < DIM; i++ {
			val, _ := strconv.ParseFloat(line[i+1], 32)
			embedding[i] = float32(val)
		}

		glove[line[0]] = embedding
	}
	return &glove
}
