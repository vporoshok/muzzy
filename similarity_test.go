package muzzy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vporoshok/muzzy"
)

func TestSimilarity(t *testing.T) {
	cases := [...]struct {
		s1, s2        string
		threshold     float64
		L, D, J, W, N float64
	}{
		{
			"happiness", "princess", 0,
			0.555, 0.555, 0.805, 0.805, 0.286,
		},
		{
			"fluffy", "fulffy", 0,
			0.666, 0.833, 0.888, 0.899, 0.5,
		},
		{
			"", "", 0,
			1, 1, 1, 1, 1,
		},
		{
			"", "any", 0,
			0, 0, 0, 0, 0,
		},
		{
			"happiness", "princess", 0.9,
			0, 0, 0, 0, 0,
		},
		{
			"abcde", "fghij", 0,
			0, 0, 0, 0, 0,
		},
		{
			"Здесь какой-то действительно большой и длинный текст с опечаткой",
			"Здесь какой-то действительно большой и длинный текст с очепаткой",
			0.9,
			0.983, 0.983, 0.989, 1.047, 0.921,
		},
		{
			"aab",
			"baa",
			0,
			0.333, 0.333, 0.777, 0.777, 0,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("%s/%s", c.s1, c.s2), func(t *testing.T) {
			assert.InDelta(t, c.L,
				muzzy.Similarity(c.s1, c.s2, muzzy.Levenshtein, c.threshold), 0.001)
			assert.InDelta(t, c.D,
				muzzy.Similarity(c.s1, c.s2, muzzy.DamerauLevenshtein, c.threshold), 0.001)
			assert.InDelta(t, c.J,
				muzzy.Similarity(c.s1, c.s2, muzzy.Jaro, c.threshold), 0.001)
			assert.InDelta(t, c.W,
				muzzy.Similarity(c.s1, c.s2, muzzy.JaroWinkler, c.threshold), 0.001)
			assert.InDelta(t, c.N,
				muzzy.Similarity(c.s1, c.s2, muzzy.NGram, c.threshold), 0.001)
		})
	}
}

//nolint:funlen // long text for test
func BenchmarkSimilarity(b *testing.B) {
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
	b.Run("Levenstein", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d := muzzy.Similarity(s1, s2, muzzy.Levenshtein, 0)
			_ = d
		}
	})
	b.Run("Damerau", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d := muzzy.Similarity(s1, s2, muzzy.DamerauLevenshtein, 0)
			_ = d
		}
	})
	b.Run("Jaro", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d := muzzy.Similarity(s1, s2, muzzy.Jaro, 0)
			_ = d
		}
	})
	b.Run("3-grams", func(b *testing.B) {
		splitter := muzzy.NGramSplitter(3, true)
		for i := 0; i < b.N; i++ {
			d := splitter.Similarity(s1, s2)
			_ = d
		}
	})
}
