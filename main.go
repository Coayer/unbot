package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/bert"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/knowledge"
	"github.com/Coayer/unbot/internal/memory"
	"github.com/Coayer/unbot/internal/plane"
	"github.com/Coayer/unbot/internal/utils"
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
	if !checkAuth(r.Header.Get("Authorization")) {
		http.Error(w, "Incorrect authentication key", http.StatusForbidden)
		log.Println("Connection denied")
		return
	}

	if r.Method == "GET" {
		query := r.URL.Query().Get("query")

		if query == "" {
			log.Println("Ping received")
		} else {
			log.Println(query)
			result := getResponse(query)
			log.Println(result)
			fmt.Fprint(w, result)
			fmt.Println()
		}
	}
}

func checkAuth(auth string) bool {
	for _, key := range utils.Config.UnbotKeys {
		if auth == key {
			return true
		}
	}
	return false
}

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? [-+x/^]")

func getResponse(query string) string {
	if strings.Contains(query, "plane") {
		return plane.GetPlane(query)
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
