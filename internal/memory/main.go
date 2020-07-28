package memory

import (
	"encoding/json"
	"github.com/Coayer/unbot/internal/bert"
	"github.com/Coayer/unbot/internal/utils"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"
)

const MEMORYPATH = "data/memories.json"

type Memory struct {
	Value  string
	Expiry int64
}

func Match(query string) bool {
	memories := readMemories()
	queryVec := sentence2Vec(utils.RemoveStopWords(query))

	for _, memory := range memories {
		similarity := cosineSimilarity(queryVec, sentence2Vec(utils.RemoveStopWords(memory.Value)))
		log.Println("Memory similarity:", similarity)
		if similarity > 0.8 {
			return true
		}
	}
	return false
}

func cosineSimilarity(a []float32, b []float32) float64 {
	var dotProduct float32
	var aSquareSum, bSquareSum float32
	for i := 0; i < DIM; i++ {
		dotProduct += a[i] * b[i]
		aSquareSum += a[i] * a[i]
		bSquareSum += b[i] * b[i]
	}
	return float64(dotProduct) / (math.Sqrt(float64(aSquareSum)) * math.Sqrt(float64(bSquareSum)))
}

func Recall(query string) string {
	log.Println("Recalling")
	memories := readMemories()
	var allMemories strings.Builder

	for _, memory := range memories {
		allMemories.WriteString(memory.Value + ". ")
	}

	return bert.AskBert(query, allMemories.String())
}

func Remember(query string) string {
	memories := readMemories()
	tokens := utils.BaseTokenize(query)[1:]

	if tokens[0] == "forever" {
		writeMemories(append(memories, Memory{Value: strings.Join(tokens[1:], " ")}))
		return "Permanent memory stored"
	} else {
		writeMemories(append(memories, Memory{Value: strings.Join(tokens, " "), Expiry: time.Now().Unix() + 86400}))
		return "Memory stored"
	}
}

func writeMemories(memories []Memory) {
	data, err := json.Marshal(memories)
	if err != nil {
		log.Fatal(err)
	}

	if ioutil.WriteFile(MEMORYPATH, data, 600) != nil {
		log.Fatal(err)
	}
}

func readMemories() []Memory {
	var memories []Memory

	data, err := ioutil.ReadFile(MEMORYPATH)
	if err != nil {
		log.Fatal(err)
	}

	if json.Unmarshal(data, &memories) != nil {
		return memories
	}

	return forget(memories)
}

func forget(memories []Memory) []Memory {
	var currentMemories []Memory
	for _, memory := range memories {
		if memory.Expiry > time.Now().Unix() || memory.Expiry == 0 {
			currentMemories = append(currentMemories, memory)
		}
	}
	return currentMemories
}
