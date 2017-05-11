package ssort_dpair

// implements qselem interface

type doublepair struct {
    x, y float64
}

func (p doublepair) Less(i interface{}) bool {
    p2 := i.(doublepair)
    return p.x < p2.x
}

type pairs []doublepair

func (s pairs) Len() int {
	return len(s)
}
func (s pairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s pairs) Less(i, j int) bool {
	return s[i].x < s[j].x
}