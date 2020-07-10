package utils

import (
	"github.com/jdkato/prose/v2"
	"log"
)

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
