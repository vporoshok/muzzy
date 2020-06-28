package muzzy_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vporoshok/muzzy"
)

type SplitIndexSuite struct {
	suite.Suite
	index *muzzy.SplitIndex
}

func (s *SplitIndexSuite) SetupSuite() {
	corpus, err := ioutil.ReadFile(filepath.Join("testdata", "dead_souls.txt"))
	s.Require().NoError(err)

	lines := strings.Split(string(corpus), "\n")
	splitter := muzzy.NGramSplitter(3, true)
	s.index = muzzy.NewSplitIndex(splitter)
	s.index.Add(lines...)
}

func (s *SplitIndexSuite) Test() {
	cases := [...]struct {
		name    string
		search  string
		nearest string
	}{
		{
			"exactly",
			`"Что ж барин? у себя, что ли?"`,
			`"Что ж барин? у себя, что ли?"`,
		},
		{
			"misstype",
			`"Что ж баирн? у себя, что ли?"`,
			`"Что ж барин? у себя, что ли?"`,
		},
		{
			"substring",
			`Что ж барин? у себя?`,
			`"Что ж барин? у себя, что ли?"`,
		},
		{
			"not found",
			"not found",
			"",
		},
	}

	// nolint:gocritic
	for _, c := range cases {
		c := c
		s.Run(c.name, func() {
			i := s.index.Search(c.search)
			s.Equal(c.nearest, s.index.Get(i))
		})
	}
}

func TestSplitIndex(t *testing.T) {
	suite.Run(t, new(SplitIndexSuite))
}
