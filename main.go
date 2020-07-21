package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/bert"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/knowledge"
	"github.com/Coayer/unbot/internal/memory"
	"github.com/Coayer/unbot/internal/plane"
	"github.com/Coayer/unbot/internal/weather"
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

BM25 stop 0.000
Config file -- location, owm key
Conversion package
*/

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Bert session")
		bert.Model.Session.Close()
		os.Exit(0)
	}()
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":2107", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.URL.Query().Get("query")

		if query == "" {
			log.Println("Ping received")
		} else {
			fmt.Println()
			log.Println(query)
			result := getResponse(query)
			log.Println(result)
			fmt.Fprint(w, result)
		}
	}
}

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? [-+x/^]")

func getResponse(query string) string {
	if strings.Contains(query, "plane") {
		return plane.GetPlane()
	} else if strings.Contains(query, "weather") || strings.Contains(query, "sunset") || strings.Contains(query, "sunrise") {
		return weather.GetWeather(query)
	} else if calculatorRegex.MatchString(query) {
		return calculator.Evaluate(query)
	} else if strings.Contains(query, "remember") {
		return memory.Remember(query)
	} else if memory.Match(query) {
		return memory.Recall(query)
	} else {
		return knowledge.AskWiki(query)
	}
}
