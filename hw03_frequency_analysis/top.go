package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

type word struct {
	text  string
	count int
}

type analyzer struct {
	words []word
}

func newAnalyzer() *analyzer {
	return &analyzer{}
}

func (a *analyzer) analyze(text string, count int) []string {
	a.parse(text)

	if len(a.words) == 0 {
		return nil
	}

	a.sort()

	return a.splice(count)
}

func (a *analyzer) parse(text string) {
	reg := regexp.MustCompile(`\S+`)

	words := map[string]int{}

	for _, word := range reg.FindAllString(text, -1) {
		words[word]++
	}

	for text, count := range words {
		a.words = append(a.words, word{text, count})
	}
}

func (a *analyzer) sort() {
	sort.Slice(a.words, func(i, j int) bool {
		if a.words[i].count == a.words[j].count {
			return a.words[i].text < a.words[j].text
		}

		return a.words[i].count > a.words[j].count
	})
}

func (a *analyzer) splice(count int) []string {
	if count > len(a.words) {
		count = len(a.words)
	}

	result := make([]string, count)

	for i := range result {
		result[i] = a.words[i].text
	}

	return result
}

func Top10(text string) []string {
	return newAnalyzer().analyze(text, 10)
}
