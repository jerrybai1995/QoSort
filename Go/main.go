package main

import (
	"fmt"
	"qosort"
)

type doublepair struct {
	x, y float64
}

type pairs []doublepair

func (s pairs) Len() int {
	return len(s)
}
func (s pairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s pairs) Less(i, j int) bool {
	if i == 1 && j == 3 {
		fmt.Println(s[i], s[j])
	}
	return s[i].x < s[j].x
}


func main() {

	qosort.Test_qsort_parallel(40)
	qosort.Test_sort(40)
}


