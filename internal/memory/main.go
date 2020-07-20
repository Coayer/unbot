package memory

import (
	"encoding/json"
	"github.com/Coayer/unbot/internal/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const DIM = 50
const MEMORYPATH = "data/memories.json"

//var GloVe = readGloVe()

type Memory struct {
	Value  string
	Expiry int64
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

//func readGloVe() *map[string][]float32 {
//	glove := make(map[string][]float32)
//	file, err := os.Open("data/glove50d.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line := strings.Split(scanner.Text(), " ")
//
//		embedding := make([]float32, DIM)
//		for i := 0; i < DIM; i++ {
//			val, _ := strconv.ParseFloat(line[i+1], 32)
//			embedding[i] = float32(val)
//		}
//
//		glove[line[0]] = embedding
//	}
//	return &glove
//}
