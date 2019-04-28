package muzzy_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vporoshok/muzzy"
)

func TestSimilarity(t *testing.T) {
	cases := [...]struct {
		s1, s2    string
		threshold float64
		L, D      float64
	}{
		{
			"", "", 0,
			1, 1,
		},
		{
			"", "any", 0,
			0, 0,
		},
		{
			"happiness", "princess", 0,
			0.555, 0.555,
		},
		{
			"happiness", "princess", 0.6,
			0, 0,
		},
		{
			"fluffy", "fulffy", 0,
			0.666, 0.833,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("%s/%s", c.s1, c.s2), func(t *testing.T) {
			assert.InDelta(t, c.L,
				muzzy.Similarity(c.s1, c.s2, muzzy.Levenshtein, c.threshold), 0.001)
			assert.InDelta(t, c.D,
				muzzy.Similarity(c.s1, c.s2, muzzy.DamerauLevenshtein, c.threshold), 0.001)
		})
	}
}
