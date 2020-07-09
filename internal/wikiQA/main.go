package wikiQA

func AskWiki(query string) string {
	articles := getArticles(query)
	best, secondBest := getRelevantArticle(articles, query)

	return askBert(query, (*articles)[best].content+" "+(*articles)[secondBest].content)
}
