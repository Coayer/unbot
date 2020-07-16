package knowledge

import (
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"hash/fnv"
	"log"
	"math"
	"sync"
)

//getRelevantArticle finds the index of the two most relevant articles
func getRelevantArticle(articles *[]article, query string) (int, int) {
	hashedQuery, hashedArticles := articleQueryBigrams(query, articles)

	idfValues := idf(hashedQuery, hashedArticles)
	avgArticleLen := meanArticleLen(hashedArticles)

	//don't need special concurrency measures as everything is fixed size
	var wait sync.WaitGroup
	wait.Add(len(*hashedArticles))

	scores := make([]float64, len(*hashedArticles))
	for i, hashedArticle := range *hashedArticles {
		go bm25(hashedQuery, hashedArticle, idfValues, avgArticleLen, i, &scores, &wait)
	}
	wait.Wait()

	for i, score := range scores {
		log.Println((*articles)[i].title + ": " + fmt.Sprintf("%f", score))
	}

	bestIndex, secondBestIndex := bestTwoScores(scores)
	return bestIndex, secondBestIndex
}

//bm25 calculates the relevance of an article to a query from within a GoRoutine
func bm25(query []int, article []int, idfValues []float64, avgDocLength float64, scoresIndex int, scores *[]float64, wait *sync.WaitGroup) {
	k := 1.2
	b := 0.75

	summands := make([]float64, len(idfValues))
	repeatedDenominatorPart := 1 - b + b*float64(len(article))/avgDocLength //saves computation as it is repeated

	for i, qi := range query {
		tf := float64(termFreq(qi, article))
		summands[i] = idfValues[i] * (tf * (k + 1) / (tf + k*repeatedDenominatorPart))
	}
	(*scores)[scoresIndex] = sum(summands)
	//+ math.Pow(math.E, 1-0.1*float64(scoresIndex))
	//experiment to use wiki relevance ranking to inform score, but order from wikipedia is page number
	wait.Done()
}

//bestTwoScores determines the indices of the highest two values of a list
func bestTwoScores(scores []float64) (int, int) {
	var bestIndex, secondBestIndex int
	highScore := math.Inf(-1)
	secondScore := math.Inf(-1)

	for i, score := range scores {
		if score > highScore {
			secondBestIndex = bestIndex
			secondScore = highScore
			bestIndex = i
			highScore = score
		} else if score > secondScore {
			secondBestIndex = i
			secondScore = score
		}
	}

	return bestIndex, secondBestIndex
}

//meanArticleLen calculates the mean of the article lengths
func meanArticleLen(articles *[][]int) float64 {
	articleLens := make([]float64, len(*articles))
	for i, article := range *articles {
		articleLens[i] = float64(len(article))
	}

	return sum(articleLens) / float64(len(*articles))
}

//idf calculates inverse document frequencies of each word in a query
func idf(query []int, documents *[][]int) []float64 {
	queryTermIDFs := make([]float64, len(query))
	for i, qi := range query {
		nQi := docsWithTerm(qi, documents)
		queryTermIDFs[i] = math.Log((float64(len(*documents)-nQi) + 0.5) / (float64(nQi) + 0.5))
	}

	return queryTermIDFs
}

//docsWithTerm calculates the number of documents containing a word
func docsWithTerm(term int, documents *[][]int) int {
	hasTerm := 0
	for _, document := range *documents {
		if termFreq(term, document) > 0 {
			hasTerm++
		}
	}

	return hasTerm
}

//termFreq calculates the times a term appears in a document
func termFreq(term int, document []int) int {
	frequency := 0
	for _, word := range document {
		if term == word {
			frequency++
		}
	}

	return frequency
}

//articleQueryBigrams produces bigram representations of a query and articles
func articleQueryBigrams(query string, articles *[]article) ([]int, *[][]int) {
	queryTokens := utils.BaseTokenize(query)
	queryBigrams := hashBigrams(queryTokens)

	articleBigrams := make([][]int, len(*articles))
	for i, article := range *articles {
		articleBigrams[i] = hashBigrams(article.tokens)
	}

	return queryBigrams, &articleBigrams
}

//converts a set of tokens into hashed bigrams
func hashBigrams(tokens []string) []int {
	ngrams := makeBigrams(tokens)
	temp := make([]int, len(ngrams))

	for i := 0; i < len(ngrams); i++ {
		hash := fnv.New32a()
		hash.Write([]byte((ngrams)[i]))
		temp[i] = int(hash.Sum32()) % 16777216 //2^24 (as used in DrQA doc retrieval system)
	}

	return temp
}

//produces bigrams for a set of tokens
func makeBigrams(tokens []string) []string {
	n := len(tokens)

	ngrams := make([]string, n*2-1)
	ngrams[0] = tokens[0]
	for i := 1; i < n; i++ {
		ngrams[i] = tokens[i]
		ngrams[i+n-1] = tokens[i-1] + tokens[i]
	}

	return ngrams
}

//calculates a sum of a set of floats
func sum(values []float64) float64 {
	result := 0.0
	for _, val := range values {
		result += val
	}

	return result
}
