package main

import (
	"fmt"
	"github.com/Coayer/unbot/internal/wikiQA"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
		fmt.Fprint(w, r.URL.Query().Get("query"))
	}
}
