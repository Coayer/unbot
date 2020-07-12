package wikiQA

import (
	"encoding/json"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"strings"
)

const wikipediaBaseURL string = "https://en.wikipedia.org/w/api.php?action=query&format=json&"

//article stores a Wikipedia article
type article struct {
	title   string
	content string
	tokens  []string
}

//ArticleSearch describes JSON format of Wikipedia search for articles
type ArticleSearch struct {
	Query struct {
		Search []struct {
			Title string
		}
	}
}

//ArticleFetch describes JSON format of fetching Wikipedia article
type ArticleFetch struct {
	Query struct {
		Pages []struct {
			Title   string
			Extract string
		}
	}
}

//makes a list of complete articles from a query
func getArticles(query string) *[]article {
	searchURL := constructTitleSearch(query)
	log.Println(searchURL)

	searchResults := utils.HttpGet(searchURL)
	titles := parseTitles(searchResults)
	log.Println(titles)

	articlesURL := constructArticleFetch(titles)
	log.Println(articlesURL)

	return parseArticles(utils.HttpGet(articlesURL))
}

//parses list of titles from JSON
func parseTitles(bytes []byte) []string {
	var response ArticleSearch
	var titles []string

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
		return titles
	}

	for _, article := range response.Query.Search {
		titles = append(titles, article.Title)
	}
	return titles
}

//parses articles from JSON
func parseArticles(bytes []byte) *[]article {
	var response ArticleFetch
	var articles []article

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
		return &[]article{}
	}

	for _, page := range response.Query.Pages {
		articles = append(articles, article{title: page.Title,
			content: page.Extract,
			tokens:  utils.BaseTokenize(page.Extract)})
	}

	return &articles
}

//creates API call from formatted query to search for relevant articles
func constructTitleSearch(srsearch string) string {
	srqiprofile := "popular_inclinks" //"engine_autoselect"
	srwhat := "text"
	srsort := "relevance" //"relevance", "none", "just_match"
	srlimit := "5"

	return wikipediaBaseURL + "list=search&srsearch=" + strings.ReplaceAll(srsearch, " ", "%20") +
		"&srlimit=" + srlimit + "&srqiprofile=" + srqiprofile + "&srwhat=" + srwhat +
		"&srinfo=&srprop=&srsort=" + srsort
}

//creates API call to fetch a given article
func constructArticleFetch(titles []string) string {
	var fetchTitles string
	first := true

	for _, title := range titles {
		if strings.Contains(title, " ") {
			title = strings.ReplaceAll(title, " ", "%20")
		}

		if !first {
			fetchTitles += "%7C" + title
		} else {
			fetchTitles += title
			first = false
		}
	}

	return wikipediaBaseURL + "prop=extracts&titles=" + fetchTitles + "&formatversion=2&exintro=1&explaintext=1"
}
