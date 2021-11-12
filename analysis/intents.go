package analysis

import (
	"encoding/json"
	"sort"

	"github.com/olivia-ai/olivia/modules"
	"github.com/olivia-ai/olivia/util"
)

var intents = map[string][]Intent{}

// Intent is a way to group sentences that mean the same thing and link them with a tag which
// represents what they mean, some responses that the bot can reply and a context
type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

// Document is any sentence from the intents' patterns linked with its tag
type Document struct {
	Sentence Sentence
	Tag      string
}

// CacheIntents set the given intents to the global variable intents
func CacheIntents(locale string, _intents []Intent) {
	intents[locale] = _intents
}

// GetIntents returns the cached intents
func GetIntents(locale string) []Intent {
	return intents[locale]
}

// SerializeIntents returns a list of intents retrieved from the given intents file
func SerializeIntents(locale string) (_intents []Intent) {
	err := json.Unmarshal(util.ReadFile("res/locales/"+locale+"/intents.json"), &_intents)
	if err != nil {
		panic(err)
	}

	CacheIntents(locale, _intents)

	return _intents
}

// SerializeModulesIntents retrieves all the registered modules and returns an array of Intents
func SerializeModulesIntents(locale string) []Intent {
	// 注意加载方式是通过init函数进行，类似 res/locales/en/modules.go:init()
	registeredModules := modules.GetModules(locale)
	intents := make([]Intent, len(registeredModules))

	for k, module := range registeredModules {
		intents[k] = Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		}
	}

	return intents
}

// GetIntentByTag returns an intent found by given tag and locale
func GetIntentByTag(tag, locale string) Intent {
	for _, intent := range GetIntents(locale) {
		if tag != intent.Tag {
			continue
		}

		return intent
	}

	return Intent{}
}

// Organize intents with an array of all words, an array with a representative word of each tag
// and an array of Documents which contains a word list associated with a tag
// 更加数据进行组织，生成所有的单词数组，所有的tag数组和具体每一个语句和tag组成的文档数组
func Organize(locale string) (words, classes []string, documents []Document) {
	// Append the modules intents to the intents from res/datasets/intents.json
	// 初始化数据从单个静态json文件，初始化modules 数据（go代码支持自定义替换）
	intents := append(
		SerializeIntents(locale),
		SerializeModulesIntents(locale)...,
	)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {
			// Tokenize the pattern's sentence
			// 过滤&&替换些字符
			patternSentence := Sentence{locale, pattern}
			patternSentence.arrange()

			// Add each word to response
			// 所有的单词，去重加入到words
			for _, word := range patternSentence.stem() {

				if !util.Contains(words, word) {
					words = append(words, word)
				}
			}

			// Add a new document
			// 每一个句子构成一个文档
			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})
		}

		// Add the intent tag to classes
		// 把每一个tag加入到分类当中
		classes = append(classes, intent.Tag)
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}
