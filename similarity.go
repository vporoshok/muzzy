package muzzy

import (
	"math"
)

type similarityAlgorithm int8

// Available algorithms to calculate strings similarity
const (
	Levenshtein similarityAlgorithm = iota
	DamerauLevenshtein
	Jaro
	JaroWinkler
	NGram
)

// Similarity of two strings with given algorithm
//
// Similarity always return number between 0 and 1, where 0 means that strings
// `s1` and `s2` is absolutely different (depends on algorithm, but it seems
// that `s1` and `s2` does not contain common symbols). And 1 means that
// strings are the same (Jaro-Winkler algorithm may return 1 even if strings
// are different).
func Similarity(s1, s2 string, algo similarityAlgorithm, threshold float64) float64 {
	if s1 == "" {
		if s2 == "" {
			return 1
		}
		return 0
	}

	var d float64
	switch algo {
	case Levenshtein, DamerauLevenshtein:
		max := math.Max(float64(len(s1)), float64(len(s2)))
		bound := int(math.Floor((1 - threshold) * max))
		var distance int
		if algo == Levenshtein {
			distance = LevenshteinDistance(s1, s2, bound)
		} else {
			distance = DamerauDistance(s1, s2, bound)
		}
		if distance < 0 {
			return 0
		}
		d = 1 - float64(distance)/max

	case Jaro:
		d = JaroSimilarity(s1, s2)

	case JaroWinkler:
		d = JaroWinklerSimilarity(s1, s2)

	default:
		d = NGramSplitter(3, true).Similarity(s1, s2)
	}
	if d < threshold {
		return 0
	}
	return d
}
