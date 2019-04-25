package muzzy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vporoshok/muzzy"
)

func TestLevenshtein(t *testing.T) {
	cases := [...]struct {
		a, b string
		res  int
	}{
		{"Something", "Smothing", 2},
		{"Something", "Smoething", 2},
		{"Something", "Some", 5},
		{"Something", "Som", -1},
		{"Something", "Smoke the king", -1},
		{"happiness", "princess", 4},
	}

	for _, c := range cases {
		assert.Equal(t, c.res, muzzy.LevenshteinDistance(c.a, c.b, 5), "%s/%s", c.a, c.b)
	}
}
