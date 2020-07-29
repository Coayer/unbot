package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func FormatHTTPQuery(query string) string {
	return strings.ReplaceAll(query, " ", "%20")
}

//performs GET request on given URL
func HttpGet(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "unbot")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
