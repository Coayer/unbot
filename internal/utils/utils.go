package utils

import (
	"github.com/jdkato/prose/v2"
	"io/ioutil"
	"log"
	"net/http"
)

const LAT, LON = 51.5074, 0.1278

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

//performs GET request on given URL
func HttpGet(url string) []byte {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return body
}
