package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/calculator"
	"github.com/Coayer/unbot/internal/conversion"
	"github.com/Coayer/unbot/internal/pkg"
	"github.com/Coayer/unbot/internal/weather"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

/*
TODO

Change client to PUT, have the server store last response, retrieve with GET
Look into session timeout for button press usage
*/

var queries = make(map[string]chan string)
var response = make(map[string]chan string)
var done = make(map[string]bool)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Closing Bert session")
		//bert.Model.Session.Close()
		os.Exit(0)
	}()
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":2107", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	validAuth, key := checkAuth(r.Header.Get("Authorization"))

	if !validAuth {
		http.Error(w, "Incorrect authentication key", http.StatusForbidden)
		log.Println("Connection denied")
		return
	}

	switch r.Method {
	case http.MethodTrace:
		log.Println("TRACE received from " + key)
		queries[key] = make(chan string, 1)
		response[key] = make(chan string, 1)
		done[key] = false
		go session(key)
	case http.MethodGet:
		log.Println("GET received from " + key)

		query := r.URL.Query().Get("query")
		fmt.Println(query)

		if tokens := pkg.BaseTokenize(query); tokens[len(tokens)-1] == pkg.Config.EndWord {
			done[key] = true
			log.Println("Sending session end")
			w.WriteHeader(http.StatusNoContent)
		} else {
			queries[key] <- query
			fmt.Fprint(w, <-response[key])
		}
	default:
		http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
	}
}

func session(key string) {
	var query, previousQuery, result string

	for {
		if done[key] {
			fmt.Println("Closing session done")
			return
		}

		select {
		case query = <-queries[key]:
			query = pkg.RemoveStopWords(query)
		case <-time.After(5 * time.Second):
			fmt.Println("Closing session timeout")
			return
		}
		//query = pkg.RemoveStopWords(<-queries[key])
		fmt.Println(query)

		if query != previousQuery && query != "" {
			log.Println(query)
			result = getResponse(query)
			log.Println(result)
			response[key] <- result
			previousQuery = query
		} else {
			response[key] <- "x"
		}
	}
}

func checkAuth(auth string) (bool, string) {
	for _, key := range pkg.Config.UnbotKeys {
		if auth == key {
			return true, auth
		}
	}
	return false, ""
}

var calculatorRegex = regexp.MustCompile("\\d+(\\.\\d+)? [-+x/^]")
var conversionRegex = regexp.MustCompile("\\d+(\\.\\d+)?")

func getResponse(query string) string {
	//if strings.Contains(query, "plane") {
	//	return plane.GetPlane(query)
	//} else

	if strings.Contains(query, "weather") || strings.Contains(query, "sunset") || strings.Contains(query, "sunrise") {
		return weather.GetWeather(query)
	} else if calculatorRegex.MatchString(query) {
		return calculator.Evaluate(query)
	} else if conversionRegex.MatchString(query) {
		return conversion.Convert(query)
	} else {
		return "hi"
	}

	//else if strings.Contains(query, "remember") {
	//	return memory.Remember(query)
	//} else if memory.Match(query) {
	//	return memory.Recall(query)
	//} else {
	//	return knowledge.AskWiki(query)
	//}
}
