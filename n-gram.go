package muzzy

import (
	"math"
	"strings"
)

// Splitter divide string to n-grams
type Splitter interface {
	Split(string) []string
	Similarity(a, b string) float64
}

// SplitterFunc is a splitter over the split function
type SplitterFunc func(string) []string

// Split string to n-grams
func (fn SplitterFunc) Split(s string) []string {
	grams := fn(s)
	set := map[string]struct{}{}
	for _, gram := range grams {
		set[gram] = struct{}{}
	}
	var res []string
	for gram := range set {
		res = append(res, gram)
	}
	return res
}

// Similarity calculate Otsuka-Ochiai coefficient
func (fn SplitterFunc) Similarity(a, b string) float64 {
	agrams := fn.Split(a)
	bgrams := fn.Split(b)
	set := map[string]struct{}{}
	intersection := float64(0)
	for _, agram := range agrams {
		set[agram] = struct{}{}
	}
	for _, bgram := range bgrams {
		if _, ok := set[bgram]; ok {
			intersection++
		}
	}

	return intersection / math.Sqrt(float64(len(agrams)*len(bgrams)))
}

// NGramSplitter is a simple n-gram splitter
func NGramSplitter(n int, withPadding bool) Splitter {
	return SplitterFunc(func(s string) []string {
		if withPadding {
			padding := strings.Repeat(" ", n-1)
			s = padding + s + padding
		}
		runes := []rune(s)
		res := make([]string, len(runes)-n+1)
		for i := 0; i <= len(runes)-n; i++ {
			res[i] = string(runes[i : i+n])
		}

		return res
	})
}

// SplitIndex index to search string in indexed strings with n-grams
type SplitIndex struct {
	Splitter
	index   map[string][]int
	strings []string
}

// NewSplitIndex is a constructor
func NewSplitIndex(splitter Splitter) *SplitIndex {

	return &SplitIndex{
		Splitter: splitter,
		index:    map[string][]int{},
	}
}

// Add string to index
func (index *SplitIndex) Add(ss ...string) {
	n := len(index.strings)
	index.strings = append(index.strings, ss...)
	for i, s := range ss {
		ngrams := index.Split(s)
		k := n + i
		for _, ngram := range ngrams {
			index.index[ngram] = append(index.index[ngram], k)
		}
	}
}

// Get string by index
func (index *SplitIndex) Get(i int) string {
	if i < 0 || i >= len(index.strings) {
		return ""
	}

	return index.strings[i]
}

// Search index of maximal similar string in index
//
// Return -1 if no string n-gram found in index.
func (index *SplitIndex) Search(s string) int {
	for i := range index.strings {
		if index.strings[i] == s {
			return i
		}
	}
	ngrams := index.Split(s)
	counters := map[int]int{}
	for _, ngram := range ngrams {
		for _, i := range index.index[ngram] {
			counters[i]++
		}
	}
	maxIndex, maxCount := -1, -1
	for i, count := range counters {
		if count > maxCount {
			maxIndex, maxCount = i, count
		}
	}

	return maxIndex
}
