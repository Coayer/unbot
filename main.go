package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/plane"
	"github.com/Coayer/unbot/internal/weather"
	"strings"

	//"github.com/Coayer/unbot/internal/wikiQA"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

/*
TODO

BM25 tuning
Weather package
Conversion package
Plane spotting package
*/

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Bert model session")
		//wikiQA.Model.Session.Close()
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

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? (-|\\+|x|\\/|\\^)")

func getResponse(query string) string {
	if strings.Contains(query, "plane") {
		return plane.GetPlane()
	} else if calculatorRegex.MatchString(query) {
		return calculator.Evaluate(query)
	} else if strings.Contains(query, "weather") {
		return weather.GetWeather(query)
	} else {
		return "ming" //wikiQA.AskWiki(query)
	}
}
