package main

import (
	"fmt"
	"qosort"
	"runtime"
	"time"
	"sorts"
	"sort"
	"math/rand"
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

func Test_qsort_parallel_ref(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 100000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	sorts.Quicksort(pairs(A))
	fmt.Println("********** Result for Serial Quicksort **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}


func main() {

	//qosort.Test_qsort_serial(10)
	//qosort.Test_qsort_by3(10)
	//qosort.Test_qsort_parallel(10)
	//qosort.Test_sort(10)
	//
	//qosort.Test_qsort_serial(20)
	//qosort.Test_qsort_by3(20)
	//qosort.Test_qsort_parallel(20)
	//qosort.Test_sort(20)
	//
	//qosort.Test_qsort_serial(30)
	//qosort.Test_qsort_by3(30)
	//qosort.Test_qsort_parallel(30)
	//qosort.Test_sort(30)
    //
	//qosort.Test_qsort_serial(40)
	//qosort.Test_qsort_by3(40)
	qosort.Test_qsort_parallel(8)
	//Test_qsort_parallel_ref(40)
	//qosort.Test_sort(8)
}


