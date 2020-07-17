package knowledge

import (
	"github.com/Coayer/unbot/internal/utils"
	"log"
)

func AskWiki(query string) string {
	ddg := getDuckDuckGo(query)
	if ddg != "" {
		log.Println("Using DuckDuckGo")
		return askBert(query, ddg)
	} else {
		log.Println("Using Wikipedia")
		cleanQuery := utils.RemoveStopWords(query)
		articles := getArticles(cleanQuery)
		best, secondBest := getRelevantArticle(articles, cleanQuery)
		return askBert(query, (*articles)[best].content+" "+(*articles)[secondBest].content)
	}
}
