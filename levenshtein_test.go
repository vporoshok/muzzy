package muzzy_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"

	"github.com/vporoshok/muzzy"
)

func TestLevenshteinDistance(t *testing.T) {
	cases := [...]struct {
		a, b     string
		max, res int
	}{
		{"Something", "Smothing", 2, 2},
		{"Something", "Smoething", 2, 2},
		{"Something", "Some", 5, 5},
		{"Something", "Som", 5, -1},
		{"Something", "Smoke the king", 7, 6},
		{"happiness", "princess", 4, 4},
		{"accabb", "bbabbabb", 4, 4},
		{"abba", "abba", 0, 0},
		{"abba", "abbb", 0, -1},
	}

	for _, c := range cases {
		assert.Equal(t, c.res, muzzy.LevenshteinDistance(c.a, c.b, c.max), "%s/%s", c.a, c.b)
	}
}

func TestDamerauDistance(t *testing.T) {
	cases := [...]struct {
		a, b     string
		max, res int
	}{
		{"Something", "Smothing", 2, 2},
		{"Something", "Smoething", 2, 1},
		{"Something", "Some", 5, 5},
		{"Something", "Som", 5, -1},
		{"Something", "Smoke the king", 6, 6},
		{"happiness", "princess", 4, 4},
		{"accabb", "bbabbabb", 4, 4},
		{"abba", "abba", 0, 0},
		{"abba", "abbb", 0, -1},
		{"abba", "baab", 2, 2},
	}

	for _, c := range cases {
		assert.Equal(t, c.res, muzzy.DamerauDistance(c.a, c.b, c.max), "%s/%s", c.a, c.b)
	}
}

func BenchmarkDistances(b *testing.B) {
	join := func(chunks ...string) string { return strings.Join(chunks, " ") }
	s1 := join(
		"В ворота гостиницы губернского города NN въехала довольно красивая",
		"рессорная небольшая бричка, в какой ездят холостяки: отставные",
		"подполковники, штабс-капитаны, помещики, имеющие около сотни душ",
		"крестьян, - словом, все те, которых называют господами средней руки.",
		"В бричке сидел господин, не красавец, но и не дурной наружности, ни",
		"слишком толст, ни слишком тонок; нельзя сказать, чтобы стар, однако ж",
		"и не так, чтобы слишком молод. Въезд его не произвел в городе",
		"совершенно никакого шума и не был сопровожден ничем особенным; только",
		"два русские мужика, стоявшие у дверей кабака против гостиницы,",
		"сделали кое-какие замечания, относившиеся, впрочем, более к экипажу,",
		"чем к сидевшему в нем. «Вишь ты, - сказал один другому, - вон какое",
		"колесо! что ты думаешь, доедет то колесо, если б случилось, в Москву",
		"или не доедет?» - «Доедет», - отвечал другой. «А в Казань-то, я",
		"думаю, не доедет?» - «В Казань не доедет», - отвечал другой. Этим",
		"разговор и кончился Да еще, когда бричка подъехала к гостинице,",
		"встретился молодой человек в белых канифасовых панталонах, весьма",
		"узких и коротких, во фраке с покушеньями на моду, из-под которого",
		"видна была манишка, застегнутая тульскою булавкою с бронзовым",
		"пистолетом. Молодой человек оборотился назад, посмотрел экипаж,",
		"придержал рукою картуз, чуть не слетевший от ветра, и пошел своей",
		"дорогой.",
	)
	s2 := join(
		"В ворота гостиницы губернского городка NN въехала довольно красивая",
		"рессорная небольшая бричка, в какой ездят холостяки: отставные",
		"подполковники, штаб-капитаны, помещики, имеющие около сотни душ",
		"крестьян, - словом, все те, которых называют господами средней руки.",
		"В бричке сидел господин, не красавец, но и не дурной наружности, ни",
		"слишком толст, ни слишком тонок; нельзя сказать, чтобы стар, однако же",
		"и не так, чтобы слишком молод. Въезд его не произвел в городе",
		"совершенно никакого шума и не был сопровожден ничем особенным; только",
		"два русские мужика, стоявшие у дверей кабака против гостиницы,",
		"сделали кое-какие замечания, относившиеся, впрочем, более к экипажу,",
		"чем к сидевшему в нем. «Видишь ты, - сказал один другому, - вон какое",
		"колесо! что ты думаешь, доедет то колесо, если б случилось, в Москву",
		"или не доедет?» - «Доедет», - отвечал другой. «А в Казань-то, я",
		"думаю, не доедет?» - «В Казань не доедет», - отвечал другой. Этим",
		"разговор и кончился Да еще, когда бричка подъехала к гостинице,",
		"встретился молодой человек в белых канифасовых панталонах, весьма",
		"узких и коротких, во фраке с замашками на моду, из-под которого",
		"видна была манишка, застегнутая тульскою булавкою с бронзовым",
		"пистолетом. Молодой человек оборотился назад, посмотрел экипаж,",
		"придержал рукою картуз, чуть не слетевший от ветра, и пошел своей",
		"дорогой.",
	)
	b.ReportAllocs()
	fmt.Printf(
		"Calculating distances between |s1|=%d and |s2|=%d (%d)\n",
		len(s1), len(s2), muzzy.LevenshteinDistance(s1, s2, -1),
	)
	bounds := [...]int{10, 15, 20, 100, 200, 500, 1000, -1}
	for _, bound := range bounds {
		bound := bound
		b.Run("Levenstein "+strconv.Itoa(bound), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d := muzzy.LevenshteinDistance(s1, s2, bound)
				_ = d
			}
		})
		b.Run("Damerau "+strconv.Itoa(bound), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d := muzzy.DamerauDistance(s1, s2, bound)
				_ = d
			}
		})
	}
}

func TestDistanceProperties(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	properties := gopter.NewProperties(nil)

	properties.Property("Levenshtein distance less or equal to changes", prop.ForAll(
		func(pair Pair) bool {
			d := muzzy.LevenshteinDistance(string(pair.a), string(pair.b), -1)
			if d > pair.changes {
				t.Logf("%s d=%d", pair, d)
			}
			return d <= pair.changes
		},
		PairGenerator(),
	))
	properties.Property("Levenshtein bounded distance same as unbounded", prop.ForAll(
		func(pair Pair) bool {
			bo := muzzy.LevenshteinDistance(string(pair.a), string(pair.b), pair.changes)
			un := muzzy.LevenshteinDistance(string(pair.a), string(pair.b), -1)
			if bo != un {
				t.Logf("%s bo=%d, un=%d", pair, bo, un)
			}
			return bo == un
		},
		PairGenerator(),
	))
	properties.Property("Damerau–Levenshtein distance less or Levenshtein", prop.ForAll(
		func(pair Pair) bool {
			l := muzzy.LevenshteinDistance(string(pair.a), string(pair.b), pair.changes)
			d := muzzy.DamerauDistance(string(pair.a), string(pair.b), pair.changes)
			if d > l {
				t.Logf("%s %d > %d", pair, d, l)
			}
			return d <= l
		},
		PairGenerator(),
	))

	properties.TestingRun(t)
}

type Pair struct {
	a, b    []rune
	changes int
}

func (pair Pair) String() string {

	return fmt.Sprintf("%s / %s [%d]", string(pair.a), string(pair.b), pair.changes)
}

func PairShrinker(v interface{}) gopter.Shrink {
	pair := v.(Pair)
	length := len(pair.a)
	if length > len(pair.b) {
		length = len(pair.b)
	}
	offset := 0
	chunk := length >> 1

	return func() (interface{}, bool) {
		if offset >= length {
			offset = 0
			chunk >>= 1
		}
		if chunk == 0 {
			return nil, false
		}
		next := Pair{
			pair.a[offset : offset+chunk],
			pair.b[offset : offset+chunk],
			pair.changes,
		}
		offset += chunk

		return next, true
	}
}

func PairGenerator() gopter.Gen {

	return func(params *gopter.GenParameters) *gopter.GenResult {
		a := gen.SliceOf(gen.Rune())(params).Result.([]rune)
		n := 10
		if len(a) < n {
			n = len(a) / 2
		}
		b := make([]rune, len(a))
		copy(b, a)
		for i := 0; i < n; i++ {
			x := params.Rng.Intn(len(b))
			y := params.Rng.Intn(len(b))
			switch params.Rng.Intn(3) {
			case 0:
				copy(b[x:], b[x+1:])
				b = b[:len(b)-1]
			case 1:
				b[x] = b[y]
			case 2:
				b = append(b[:x], append([]rune{b[y]}, b[x:]...)...)
			}
		}

		return gopter.NewGenResult(Pair{a, b, n}, PairShrinker)
	}
}
