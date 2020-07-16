package knowledge

import (
	"encoding/json"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"strings"
	"unicode"
)

const wikipediaBaseURL string = "https://en.wikipedia.org/w/api.php?action=query&format=json&"

var stopWords = []string{"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your", "yours", "yourself",
	"yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its", "itself", "they", "them",
	"their", "theirs", "themselves", "what", "which", "who", "whom", "this", "that", "these", "those", "am", "is", "are",
	"was", "were", "be", "been", "being", "have", "has", "had", "having", "do", "does", "did", "doing", "a", "an", "the",
	"and", "but", "if", "or", "because", "as", "until", "while", "of", "at", "by", "for", "with", "about", "against",
	"between", "into", "through", "during", "before", "after", "above", "below", "to", "from", "up", "down", "in", "out",
	"on", "off", "over", "under", "again", "further", "then", "once", "here", "there", "when", "where", "why", "how", "all",
	"any", "both", "each", "few", "more", "most", "other", "some", "such", "no", "nor", "not", "only", "own", "same", "so",
	"than", "too", "very", "can", "will", "just", "don't", "should", "now"}

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
	query = removeStopWords(query)
	searchURL := constructTitleSearch(query)
	log.Println(searchURL)

	searchResults := utils.HttpGet(searchURL)
	titles := parseTitles(searchResults)
	log.Println(titles)

	articlesURL := constructArticleFetch(titles)
	log.Println(articlesURL)

	return parseArticles(utils.HttpGet(articlesURL))
}

func removeStopWords(query string) string {
	var cleanedQuery strings.Builder

	for _, token := range utils.BaseTokenize(query) {
		for i := 0; i < len(stopWords); i++ {
			if stopWords[i] == token {
				break
			}
			if i == len(stopWords)-1 {
				cleanedQuery.WriteString(token)
			}
		}
	}

	return cleanedQuery.String()
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
		if isASCII(article.Title) {
			titles = append(titles, article.Title)
		}
	}
	return titles
}

func isASCII(title string) bool {
	for i := 0; i < len(title); i++ {
		if title[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
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
