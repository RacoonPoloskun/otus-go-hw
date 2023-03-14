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

func (anl *analyzer) analyze(text string, count int) []string {
	anl.parse(text)

	if len(anl.words) == 0 {
		return nil
	}

	anl.sort()

	return anl.splice(count)
}

func (anl *analyzer) parse(text string) {
	reg, _ := regexp.Compile(`\S+`)
	words := map[string]int{}

	for _, word := range reg.FindAllString(text, -1) {
		words[word]++
	}

	for text, count := range words {
		anl.words = append(anl.words, word{text, count})
	}
}

func (anl *analyzer) sort() {
	sort.Slice(anl.words, func(i, j int) bool {
		if anl.words[i].count == anl.words[j].count {
			return anl.words[i].text < anl.words[j].text
		}

		return anl.words[i].count > anl.words[j].count
	})
}

func (anl *analyzer) splice(count int) []string {
	if count > len(anl.words) {
		count = len(anl.words)
	}

	result := make([]string, count)

	for i, _ := range result {
		result[i] = anl.words[i].text
	}

	return result
}

func Top10(text string) []string {
	return newAnalyzer().analyze(text, 10)
}
