package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/conversion"
	"github.com/Coayer/unbot/internal/knowledge"
	"github.com/Coayer/unbot/internal/memory"
	"github.com/Coayer/unbot/internal/pkg"
	"github.com/Coayer/unbot/internal/pkg/bert"
	"github.com/Coayer/unbot/internal/plane"
	"github.com/Coayer/unbot/internal/reminder"
	"github.com/Coayer/unbot/internal/weather"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

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

func handler(w http.ResponseWriter, request *http.Request) {
	validAuth := checkAuth(request.Header.Get("Authorization"))

	if !validAuth {
		http.Error(w, "Incorrect authentication key", http.StatusForbidden)
		log.Println("Connection denied")
		return
	}

	switch request.Method {
	case http.MethodHead:
		log.Println("HEAD received")
	case http.MethodGet:
		log.Println("GET received")

		var response string

		if condition := request.URL.Query().Get("reminder"); condition != "" {
			response = reminder.GetReminders(condition)
		} else {
			query := request.URL.Query().Get("query")
			response = getResponse(query)

			log.Println(query)
			log.Println(response)
		}

		_, err := fmt.Fprint(w, response)
		if err != nil {
			log.Println(err)
		}
	default:
		http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
	}
}

func checkAuth(auth string) bool {
	for _, key := range pkg.Config.UnbotKeys {
		if auth == key {
			return true
		}
	}
	return false
}

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? [-+x/^]")
var conversionRegex = regexp.MustCompile("\\d+(\\.\\d+)? .+ in")

func getResponse(query string) string {
	if strings.Contains(query, "plane") {
		return plane.GetPlane(query)
	} else if strings.Contains(query, "weather") || strings.Contains(query, "sunset") || strings.Contains(query, "sunrise") {
		return weather.GetWeather(query)
	} else if calculatorRegex.MatchString(query) {
		return calculator.Evaluate(query)
	} else if conversionRegex.MatchString(query) {
		return conversion.Convert(query)
	} else if strings.Contains(query, "remind") {
		return reminder.SetReminder(query)
	} else if strings.Contains(query, "remember") {
		return memory.Remember(query)
	} else if memory.Match(query) {
		return memory.Recall(query)
	} else {
		return knowledge.AskWiki(query)
	}
}
