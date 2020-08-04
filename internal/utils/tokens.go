package utils

import (
	"github.com/jdkato/prose/v2"
	"log"
	"strings"
)

var stopWords = []string{"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your", "yours", "yourself",
	"yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its", "itself", "they", "them",
	"their", "theirs", "themselves", "what", "which", "who", "whom", "this", "that", "these", "those", "am", "is", "are",
	"was", "were", "be", "been", "being", "have", "has", "had", "having", "do", "does", "did", "doing", "a", "an", "the",
	"and", "but", "if", "or", "because", "as", "until", "while", "of", "at", "by", "for", "with", "about", "against",
	"between", "into", "through", "during", "before", "after", "above", "below", "to", "from", "up", "down", "in", "out",
	"on", "off", "over", "under", "again", "further", "then", "once", "here", "there", "when", "where", "why", "how", "all",
	"any", "both", "each", "few", "more", "most", "other", "some", "such", "no", "nor", "not", "only", "own", "same", "so",
	"than", "too", "very", "can", "will", "just", "don't", "should", "now"}

//simple wrapper for prose tokenizer
func BaseTokenize(sequence string) []string {
	doc, err := prose.NewDocument(sequence,
		prose.WithSegmentation(false),
		prose.WithTagging(false),
		prose.WithExtraction(false))

	if err != nil {
		log.Fatal(err)
	}

	var result []string

	for _, tok := range doc.Tokens() {
		result = append(result, tok.Text)
	}

	return result
}

func RemoveStopWords(sentence string) string {
	var cleanedQuery strings.Builder

	for _, token := range BaseTokenize(sentence) {
		for i := 0; i < len(stopWords); i++ {
			if stopWords[i] == token {
				break
			}
			if i == len(stopWords)-1 {
				cleanedQuery.WriteString(token + " ")
			}
		}
	}

	return cleanedQuery.String()
}

func GetEntities(sequence string) string {
	doc, err := prose.NewDocument(sequence,
		prose.WithSegmentation(false))

	if err != nil {
		log.Fatal(err)
	}

	var entity strings.Builder

	for _, ent := range doc.Entities() {
		log.Println("Found entity: " + ent.Text)
		entity.WriteString(ent.Text + " ")
	}

	return entity.String()
}

//isNumeric checks if a token is a number
func IsNumeric(token string) bool {
	for _, char := range token {
		if (char < '0' || char > '9') && char != '.' {
			return false
		}
	}
	return true
}
