package main

import (
	"fmt"
	"sort"
	"time"
	"math/rand"
	"runtime"
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
	return s[i].x < s[j].x
}


func main() {
	n := 10000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	qosort.QuickSort(pairs(A), 0, n)
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check they array is sorted: ", sort.IsSorted(pairs(A)))
}


