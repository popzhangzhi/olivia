package analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/olivia-ai/olivia/locales"

	"github.com/olivia-ai/olivia/util"
	"github.com/tebeka/snowball"
)

// arrange checks the format of a string to normalize it, remove ignored characters
// 简单过滤和替换，看起来对非英文不生效（替换.?!符号）
func (sentence *Sentence) arrange() {
	// Remove punctuation after letters
	punctuationRegex := regexp.MustCompile(`[a-zA-Z]( )?(\.|\?|!|¿|¡)`)
	sentence.Content = punctuationRegex.ReplaceAllStringFunc(sentence.Content, func(s string) string {
		punctuation := regexp.MustCompile(`(\.|\?|!)`)
		return punctuation.ReplaceAllString(s, "")
	})

	sentence.Content = strings.ReplaceAll(sentence.Content, "-", " ")
	sentence.Content = strings.TrimSpace(sentence.Content)
}

// removeStopWords takes an arary of words, removes the stopwords and returns it
func removeStopWords(locale string, words []string) []string {
	// Don't remove stopwords for small sentences like “How are you” because it will remove all the words
	if len(words) <= 4 {
		return words
	}

	// Read the content of the stopwords file
	stopWords := string(util.ReadFile("res/locales/" + locale + "/stopwords.txt"))

	var wordsToRemove []string
	// TODO 可以用更好的方式，第一标记需要移除的元素下标。第二次遍历加入即可。为什么当前使用的第一次取出需要移除的单词。第二次进行2次循环对比后得出最终被移除stopwrod之后的数据
	// Iterate through all the stopwords
	for _, stopWord := range strings.Split(stopWords, "\n") {
		// Iterate through all the words of the given array
		for _, word := range words {
			// Continue if the word isn't a stopword
			if !strings.Contains(stopWord, word) {
				continue
			}

			wordsToRemove = append(wordsToRemove, word)
		}
	}

	return util.Difference(words, wordsToRemove)
}

// tokenize returns a list of words that have been lower-cased
func (sentence Sentence) tokenize() (tokens []string) {
	// Split the sentence in words
	tokens = strings.Fields(sentence.Content)

	// Lower case each word
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}

	tokens = removeStopWords(sentence.Locale, tokens)

	return
}

// stem returns the sentence split in stemmed words
// 过滤stopWord 并且去除词的状态
func (sentence Sentence) stem() (tokenizeWords []string) {
	locale := locales.GetNameByTag(sentence.Locale)
	// Set default locale to english
	if locale == "" {
		locale = "english"
	}
	// pattern中剔除stopWord
	// TODO 仅适配与英文等分词用空格的语言，先空格分词，然后小写单词，然后移除指定的stopword
	tokens := sentence.tokenize()

	stemmer, err := snowball.New(locale)
	if err != nil {
		fmt.Println("Stemmer error", err)
		return
	}
	// 进行分词去除一些状态，例如从running -> run
	// Get the string token and push it to tokenizeWord
	for _, tokenizeWord := range tokens {
		word := stemmer.Stem(tokenizeWord)
		tokenizeWords = append(tokenizeWords, word)
	}

	return
}

// WordsBag retrieves the intents words and returns the sentence converted in a bag of words
func (sentence Sentence) WordsBag(words []string) (bag []float64) {
	for _, word := range words {
		// Append 1 if the patternWords contains the actual word, else 0
		// 确定当前语句是否包含单词
		var valueToAppend float64
		// TODO 在之前训练数据的时候就已经进行过分词了，为什么还再次进行分词，来判断是否含有单词
		// 如果拥有词，置1，否则0
		if util.Contains(sentence.stem(), word) {
			valueToAppend = 1
		}

		bag = append(bag, valueToAppend)
	}

	return bag
}
