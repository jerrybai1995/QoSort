package main

import (
	"fmt"
	"qosort"
	"os"
	"runtime"
	"time"
	"sort"
	"math/rand"
	"strconv"
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
	cores := runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(cores)

	args := os.Args
	n := 100000000
	var err error
	if len(args) > 1 {
		n, err = strconv.Atoi(args[1])
		if err != nil { fmt.Println("Need input argument for array length n.") }
	}

	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	qosort.Qsort_parallel(pairs(A), 0, n)
	fmt.Println("********** Result for (Optimized) Parallel Quicksort **********")
	fmt.Println("Length of double-pair array: ", n)
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}


