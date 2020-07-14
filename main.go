package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/plane"
	"github.com/Coayer/unbot/internal/weather"
	"github.com/Coayer/unbot/internal/wikiQA"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

/*
TODO

Documentation
BM25 tuning
Config file -- location, owm key
Conversion package
*/

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Bert model session")
		wikiQA.Model.Session.Close()
		os.Exit(0)
	}()
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":1337", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.URL.Query().Get("query")

		if query == "" {
			log.Println("Ping received")
		} else {
			log.Println(query)
			result := getResponse(query)
			log.Println(result)
			fmt.Fprint(w, result)
		}
	}
}

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? [-+x\\/^]")

func getResponse(query string) string {
	if strings.Contains(query, "plane") {
		return plane.GetPlane()
	} else if strings.Contains(query, "weather") {
		return weather.GetWeather(query)
	} else if calculatorRegex.MatchString(query) {
		return calculator.Evaluate(query)
	} else {
		return wikiQA.AskWiki(query)
	}
}
