package muzzy

// Distance between strings
//
type Distance func(s1, s2 string, bound int) int

// LevenshteinDistance calculate distance between strings
//
// Distance between strings represent how many simple operations (insertion,
// deletion and replacement) needed to be done to convert first string to
// another. For examples, distance between "milk" and "silk" is 1 (replacement
// 'm' to 's'), but distance between "happiness" and "princess" is 4 (-'h',
// -'a', 'p'/'r', +'c').
//
// Parameter `bound` is used to optimization. If you only intrested if distance
// less or equal some number, put it number as bound. If distance more than this
// number, function return -1. Use -1 as `bound` to calculate distance without
// limitation.
func LevenshteinDistance(s1, s2 string, bound int) int {
	if bound == 0 {
		if s1 == s2 {
			return 0
		}
		return -1
	}
	b := &bounder{
		bound: bound,
	}

	return b.Do(s1, s2, newLevenshteinCalculator)
}

// DamerauDistance similar to Levenshtein except that permutation cost is 1
//
// Permutation of neighbor symbols cost is 1, for example Levenshtein distance
// between "permutation" and "permtuation" is 2 (u/t, t/u), but in
// Damerau–Levenshtein is 1.
func DamerauDistance(s1, s2 string, bound int) int {
	if bound == 0 {
		if s1 == s2 {
			return 0
		}
		return -1
	}
	b := &bounder{
		bound: bound,
	}

	return b.Do(s1, s2, newDamerauCalculator)
}

// Calculator is an abstraction of handling prefix-distance matrix to
// isolate implementation with one row (Levenshtein distance) and two rows
// (Damerau–Levenshtein distance).
type calculator interface {
	Reset(int)
	Calc(int, int) int
}

type bounder struct {
	calc   calculator
	bound  int
	width  int
	height int
	left   int
	right  int
}

func (b *bounder) Do(s1, s2 string, calc func(r1, r2 []rune) calculator) int {
	r1, r2 := []rune(s1), []rune(s2)
	if len(r1) < len(r2) {
		r1, r2 = r2, r1
	}
	if b.bound >= 0 && len(r1)-len(r2) > b.bound {
		return -1
	}
	b.width = len(r2)
	b.height = len(r1)
	if b.bound < 0 {
		b.bound = len(r1)
	}
	b.right = len(r2)
	if b.bound < b.right {
		b.right = b.bound
	}
	b.calc = calc(r1, r2)

	return b.Calculate()
}

// Calculate distance matrix
func (b *bounder) Calculate() int {
	var n int
	for i := 0; i < b.height; i++ {
		b.calc.Reset(b.left)
		for j := b.left; j < b.right; j++ {
			n = b.calc.Calc(i, j)
			if n > b.bound && j == b.left {
				b.left++
			}
		}
		if b.right < b.width && n <= b.bound {
			b.right++
		}
		if b.left >= b.right {
			return -1
		}
	}

	return n
}

type levenshteinCalculator struct {
	s1, s2 []rune
	last   []int
	diag   int
}

func newLevenshteinCalculator(s1, s2 []rune) calculator {
	lc := &levenshteinCalculator{s1: s1, s2: s2, last: make([]int, len(s2)+1)}
	for i := range lc.last {
		lc.last[i] = i
	}

	return lc
}

func (lc *levenshteinCalculator) Reset(j int) {
	lc.diag = lc.last[j]
	lc.last[j]++
}

func (lc *levenshteinCalculator) Calc(i, j int) int {
	dd := lc.diag
	if lc.s1[i] != lc.s2[j] {
		dd++
	}
	lc.last[j+1], lc.diag = min(dd, lc.last[j]+1, lc.last[j+1]+1), lc.last[j+1]

	return lc.last[j+1]
}

type damerauCalculator struct {
	s1, s2 []rune
	last   []int
	prev   []int
	buff   [2]int
}

func newDamerauCalculator(s1, s2 []rune) calculator {
	lc := &damerauCalculator{
		s1:   s1,
		s2:   s2,
		last: make([]int, len(s2)+1),
		prev: make([]int, len(s2)+1),
	}
	for i := range lc.last {
		lc.last[i] = i
		lc.prev[i] = i
	}

	return lc
}

func (lc *damerauCalculator) Reset(j int) {
	lc.rotate(j, lc.last[j]+1)
}

func (lc *damerauCalculator) Calc(i, j int) int {
	dd := lc.prev[j]
	if lc.s1[i] != lc.s2[j] {
		dd++
	}
	if i > 0 && j > 0 && lc.s1[i-1] == lc.s2[j] && lc.s1[i] == lc.s2[j-1] {

		return lc.rotate(j+1, min(dd, lc.last[j]+1, lc.last[j+1]+1, lc.buff[0]+1))
	}

	return lc.rotate(j+1, min(dd, lc.last[j]+1, lc.last[j+1]+1))
}

func (lc *damerauCalculator) rotate(j, next int) int {
	lc.buff[0], lc.buff[1], lc.prev[j], lc.last[j] =
		lc.buff[1], lc.prev[j], lc.last[j], next

	return next
}

func min(x ...int) int {
	m := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] < m {
			m = x[i]
		}
	}

	return m
}
