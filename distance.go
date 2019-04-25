package muzzy

// Distance between strings
//
// Parameter `max` is used to optimization. If you only intrested if distance
// less or equal some number, put it number as max. If distance more than this
// number, function return -1. Use -1 as `max` to calculate distance without
// limitation.
type Distance func(a, b string, max int) int

// LevenshteinDistance calculate distance between strings
//
// Distance between strings represent how many simple operations (insertion,
// deletion and replacement) needed to be done to convert first string to
// another. For examples, distance between "milk" and "silk" is 1 (replacement
// 'm' to 's'), but distance between "happiness" and "princess" is 4 (-'h',
// -'a', 'p'/'r', +'c').
//
// See Distance type for explanation of `max` param.
func LevenshteinDistance(a, b string, max int) int {
	d := &levenshteinDistance{a: []rune(a), b: []rune(b), max: max}

	return d.Calculate()
}

type levenshteinDistance struct {
	a, b []rune
	max  int
	mem  []int
	left int
}

func (d *levenshteinDistance) Calculate() int {
	if len(d.a) < len(d.b) {
		d.a, d.b = d.b, d.a
	}
	if d.max >= 0 && len(d.a)-len(d.b) > d.max {
		return -1
	}
	d.InitMem()
	for i := range d.a {
		if d.max >= 0 && d.TrimLeft() {
			return -1
		}
		diag := d.mem[d.left]
		d.mem[d.left]++
		for j := d.left; j < len(d.b); j++ {
			// if mem[j] >= k && diag > k && mem[j+1] >= k {
			// 	break
			// }
			diagDistance := diag
			if d.a[i] != d.b[j] {
				diagDistance++
			}
			d.mem[j+1], diag = min(diagDistance, d.mem[j]+1, d.mem[j+1]+1), d.mem[j+1]
		}
	}

	return d.mem[len(d.b)]
}

// InitMem allocate and initiate memory for last calculated row
func (d *levenshteinDistance) InitMem() {
	d.mem = make([]int, len(d.b)+1)
	for i := range d.mem {
		d.mem[i] = i
	}
}

func (d *levenshteinDistance) TrimLeft() bool {
	for ; d.mem[d.left] > d.max; d.left++ {
		if d.left == len(d.b) {

			return true
		}
	}

	return false
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
