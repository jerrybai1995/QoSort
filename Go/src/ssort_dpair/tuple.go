package ssort_dpair

type tuple struct {
	x int
	y int
}

type scheduler func(tuple)