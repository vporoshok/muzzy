package muzzy

// JaroSimilarity return how close s1 to s2
func JaroSimilarity(s1, s2 string) float64 {

	return newJaroCalculator(s1, s2).Do()
}

type jaroCalculator struct {
	s1, s2 []rune
	l1, l2 []bool
}

func newJaroCalculator(s1, s2 string) *jaroCalculator {
	jc := &jaroCalculator{s1: []rune(s1), s2: []rune(s2)}
	if len(jc.s1) > len(jc.s2) {
		jc.s1, jc.s2 = jc.s2, jc.s1
	}
	jc.l1 = make([]bool, len(jc.s1))
	jc.l2 = make([]bool, len(jc.s2))

	return jc
}

func (jc *jaroCalculator) Do() float64 {
	m := jc.FindMatches()
	if m == 0 {
		return 0
	}
	t := jc.FindTranspositions()
	n1, n2 := float64(len(jc.s1)), float64(len(jc.s2))

	return (m/n1 + m/n2 + (m-t)/m) / 3
}

func (jc *jaroCalculator) FindMatches() float64 {
	eps := len(jc.s2) >> 1
	m := 0.0
	for i := range jc.s1 {
		l := 0
		if i-eps > 0 {
			l = i - eps
		}
		r := len(jc.s2)
		if i+eps < len(jc.s2) {
			r = i + eps
		}
		for j := l; j < r; j++ {
			if jc.s1[i] == jc.s2[j] && !jc.l2[j] {
				jc.l1[i] = true
				jc.l2[j] = true
				m++
				break
			}
		}
	}

	return m
}

func (jc *jaroCalculator) FindTranspositions() float64 {
	t := 0.0
	i, j := 0, 0
	for ; i < len(jc.s1); i++ {
		if jc.l1[i] {
			for ; j < len(jc.s2); j++ {
				if jc.l2[j] {
					j++
					break
				}
			}
			if jc.s1[i] != jc.s2[j-1] {
				t++
			}
		}
	}

	return t
}
