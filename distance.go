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
	if max == 0 {
		if a == b {
			return 0
		}
		return -1
	}
	d := &levenshteinDistance{a: []rune(a), b: []rune(b), max: max}

	return d.Do()
}

type levenshteinDistance struct {
	a, b []rune
	max  int
	mem  []int
	l, r int
}

func (d *levenshteinDistance) Do() int {
	if len(d.a) < len(d.b) {
		d.a, d.b = d.b, d.a
	}
	if d.max >= 0 && len(d.a)-len(d.b) > d.max {
		return -1
	}
	d.Init()

	return d.Calculate()
}

// Init allocate and initiate memory for last calculated row and setup right bound
func (d *levenshteinDistance) Init() {
	d.mem = make([]int, len(d.b)+1)
	for i := range d.mem {
		d.mem[i] = i
	}
	if d.max < 0 {
		d.max = len(d.a)
	}
	d.r = len(d.b)
	if d.max < d.r {
		d.r = d.max
	}
}

// Calculate distance matrix
func (d *levenshteinDistance) Calculate() int {
	for i := range d.a {
		diag := d.mem[d.l]
		d.mem[d.l]++
		for j := d.l; j < d.r; j++ {
			diagDistance := diag
			if d.a[i] != d.b[j] {
				diagDistance++
			}
			d.mem[j+1], diag = min(diagDistance, d.mem[j]+1, d.mem[j+1]+1), d.mem[j+1]
			if d.mem[j+1] > d.max && j == d.l {
				d.l++
			}
		}
		if d.r < len(d.b) && d.mem[d.r] <= d.max {
			d.r++
		}
		if d.l >= d.r {
			return -1
		}
	}

	return d.mem[len(d.b)]
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
