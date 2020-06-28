package muzzy

import "sync"

// WinklerScalingFactor how much the score is adjusted upwards for having common prefixes.
const WinklerScalingFactor = 0.1

// JaroWinklerSimilarity return how close s1 and s2 increase similarity of same prefixed.
func JaroWinklerSimilarity(s1, s2 string) float64 {
	s := JaroSimilarity(s1, s2)
	l, r1, r2 := 0, []rune(s1), []rune(s2)
	n := len(r1)

	if len(r2) < n {
		n = len(r2)
	}

	for ; l < n && r1[l] == r2[l]; l++ {
	}

	return s + float64(l)*(1-s)*WinklerScalingFactor
}

// JaroSimilarity return how close s1 to s2
//
// To find the Jaro similarity we should find the number of matched characters $m$.
// If $m$ is equal zero then strings absolutely different and similarity is equal zero too.
// To simplify algorithm we assume that length of `s2` great or equal length of `s1`.
// If it is not, we swap strings.
//
// Characters of strings `s1` and `s2` are matched if they equal and their positions are closer
// than half of `s2` length. Of course, one character may be used only in one matched pair.
//
// If $m$ is not zero, we also should find the number of transpositions $t$,
// i.e. the number of matched pairs that switch their order.
//
// The Jaro similarity is
// \[
//   \frac{1}{3}\left(\frac{m}{\|s_1\|} + \frac{m}{\|s_2\|} + \frac{m - t}{m}\right)
// \].
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

// To find matched characters we use a cartesian tree of characters of `s2`,
// where the priority is a character position and the key is a character code.
// We will control the epsilon neighborhood of the current character in `s1`
// by adding the next character of `s2` in the tree on every step,
// and poping root if its priority too far from the current character position.
func (jc *jaroCalculator) Do() float64 {
	m := jc.FindMatchesCartesian()
	if m == 0 {
		return 0
	}

	t := jc.FindTranspositions()
	n1, n2 := float64(len(jc.s1)), float64(len(jc.s2))

	return (m/n1 + m/n2 + (m-t)/m) / 3
}

func (jc *jaroCalculator) FindMatchesCartesian() float64 {
	eps := len(jc.s2) >> 1
	tree := newCartesianTree()
	m := 0.0

	for i := 0; i < eps; i++ {
		tree.Add(jc.s2[i])
	}

	for i := range jc.s1 {
		if i+eps < len(jc.s2) {
			tree.Add(jc.s2[i+eps])
		}

		if !tree.Empty() && tree.Peek() < i-eps {
			tree.Pop()
		}

		j := tree.SearchAndDelete(jc.s1[i])
		if j != -1 {
			jc.l1[i] = true
			jc.l2[j] = true
			m++
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

type node struct {
	priority int
	key      rune
	left     *node
	right    *node
}

func (n *node) Add(m *node) {
	for {
		next := &n.right
		if m.key < n.key {
			next = &n.left
		}

		if *next == nil {
			*next = m
			return
		}

		n = *next
	}
}

func (n *node) Search(key rune) (res, parent *node) {
	res = n
	for res != nil {
		if key == res.key {
			return
		}

		parent = res

		if key < res.key {
			res = res.left
		} else {
			res = res.right
		}
	}

	return nil, nil
}

type cartesianTree struct {
	root *node
	pool *sync.Pool
	next int
}

func newCartesianTree() *cartesianTree {
	return &cartesianTree{
		pool: &sync.Pool{
			New: func() interface{} { return &node{} },
		},
	}
}

func (tree *cartesianTree) Add(key rune) {
	n := tree.pool.Get().(*node)

	n.priority = tree.next
	n.key = key
	n.left = nil
	n.right = nil

	if tree.root == nil {
		tree.root = n
	} else {
		tree.root.Add(n)
	}

	tree.next++
}

func (tree *cartesianTree) SearchAndDelete(key rune) int {
	n, p := tree.root.Search(key)
	if n == nil {
		return -1
	}

	m := tree.Merge(n.left, n.right)

	switch {
	case p == nil:
		tree.root = m
	case p.left == n:
		p.left = m
	default:
		p.right = m
	}

	res := n.priority
	tree.pool.Put(n)

	return res
}

func (tree *cartesianTree) Empty() bool {
	return tree.root == nil
}

func (tree *cartesianTree) Peek() int {
	return tree.root.priority
}

func (tree *cartesianTree) Pop() {
	m := tree.Merge(tree.root.left, tree.root.right)
	if m == nil {
		tree.root = nil
	} else {
		*tree.root = *m
		tree.pool.Put(m)
	}
}

func (tree *cartesianTree) Merge(n, m *node) *node {
	if n == nil {
		return m
	}

	if m == nil {
		return n
	}

	if m.priority < n.priority {
		n, m = m, n
	}

	if m.key < n.key {
		n.left = tree.Merge(n.left, m)
	} else {
		n.right = tree.Merge(n.right, m)
	}

	return n
}
