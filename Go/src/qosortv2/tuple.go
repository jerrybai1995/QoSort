package qosortv2

type tuple struct {
	x int
	y int
}

type scheduler func(tuple)