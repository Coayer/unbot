package knowledge

import (
	"github.com/Coayer/unbot/internal/pkg"
	"github.com/Coayer/unbot/internal/pkg/bert"
	"github.com/jdkato/prose/v2"
	"log"
	"strings"
)

var previousQuery string

func AskWiki(query string) string {
	query = addMissingEntities(query)
	previousQuery = query
	log.Println(query)
	ddg := getDuckDuckGo(query)

	if ddg != "" {
		log.Println("Using DuckDuckGo")
		return bert.AskBert(query, ddg)
	} else {
		log.Println("Using Wikipedia")

		cleanQuery := pkg.RemoveStopWords(query)
		articles := getArticles(cleanQuery)
		best, secondBest := getRelevantArticle(articles, cleanQuery)
		return bert.AskBert(query, (*articles)[best].content+" "+(*articles)[secondBest].content)
	}
}

func addMissingEntities(query string) string {
	doc, err := prose.NewDocument(query,
		prose.WithSegmentation(false),
		prose.WithExtraction(false))

	if err != nil {
		log.Fatal(err)
	}

	var newQuery strings.Builder

	for _, token := range doc.Tokens() {
		if token.Tag == "PRP" || token.Tag == "PRP$" {
			newQuery.WriteString(pkg.GetEntities(previousQuery))
		} else {
			if len(token.Text) != 1 {
				newQuery.WriteString(token.Text + " ")
			} else {
				newQuery.WriteString(token.Text)
			}
		}
	}
	newQuery.WriteString("?") //needed for prose NER to work
	return newQuery.String()
}
