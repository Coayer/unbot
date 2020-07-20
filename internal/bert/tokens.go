package bert

import (
	"bufio"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"os"
	"strings"
)

var vocab *vocabulary = loadVocab()

//vocabulary stores the tokenizer's vocabulary
type vocabulary struct {
	encodeMap map[string]int32
	vocabMap  map[string]bool
}

//loadVocab loads BERT vocabulary from text file
func loadVocab() *vocabulary {
	file, err := os.Open("data/vocab.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	encode := make(map[string]int32, 30522) //BERT vocab length
	check := make(map[string]bool, 30522)

	var i int32
	for scanner.Scan() {
		encode[scanner.Text()] = i
		check[scanner.Text()] = true
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &vocabulary{encodeMap: encode, vocabMap: check}
}

func tokensToEnglish(tokens []string) string {
	var result strings.Builder

	for _, token := range tokens {
		if len(token) != 1 && token[0:2] == "##" {
			result.WriteString(token[2:])
		} else {
			result.WriteString(" " + token)
		}
	}

	return result.String()
}

//tokenize does full wordpiece tokenization on a sequence
func tokenize(sequence string) []string {
	var result []string

	for _, token := range utils.BaseTokenize(sequence) {
		for _, piece := range wordPieceTokenizer(token) {
			result = append(result, piece)
		}
	}

	return result
}

//splits a single token into pieces if it contains subwords from vocabulary
func wordPieceTokenizer(token string) []string {
	token = strings.ToLower(token)

	//can't get pieces from single char words, need special consideration
	if len(token) == 1 {
		//checks if token is alphanumeric or specific punctuation which doesn't need ## prefix
		if (token >= "a" && token <= "z") || (token >= "0" && token <= "9") ||
			(token == "(" || token == "[" || token == "{" || token == "~") {
			return []string{token}
		}
		return []string{"##" + token} //hard coding rules is not ideal, but will work most of the time
	}

	var pieces []string
	lastSplit := 0
	firstSubWord := true

	/*goes forwards along letter indices, and backwards from end of word to current letter
	this means it can match longest pieces first by check for matches with vocabulary*/
	for i := 1; i <= len(token); i++ {
		for z := len(token); z >= i; z-- {
			//won't place ## in front of piece if at start of word
			if firstSubWord {
				piece := token[lastSplit:z]
				if vocab.vocabMap[piece] {
					pieces = append(pieces, piece)
					lastSplit = i
					i = z //jumps ahead as there is a piece now "missing" from the word
					firstSubWord = false
					break
				}
			} else {
				piece := "##" + token[i-1:z]
				if vocab.vocabMap[piece] {
					pieces = append(pieces, piece)
					lastSplit = z
					i = z
					break
				}
			}
		}
	}

	return pieces
}
