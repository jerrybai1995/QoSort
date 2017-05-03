package qosort

import (
	"fmt"
	"sort"
	"time"
	"math/rand"
	"runtime"
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

func Test_qsort_parallel(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 100000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	Qsort_parallel(pairs(A), 0, n)
	fmt.Println("********** Result for Serial Quicksort **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}

func Test_qsort_serial(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 100000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	Qsort_serial(pairs(A), 0, len(A))
	fmt.Println("********** Result for Serial Quicksort **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}

func Test_qsort_by3(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 100000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	qsort_by3(pairs(A), 0, len(A))
	fmt.Println("********** Result for Parallel Quicksort (Split by 3) **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}

func Test_qsort_par(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 1000000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	QuickSort(pairs(A))
	fmt.Println("********** Result for Parallel Quicksort (Split by 2) **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}

func Test_sort(cores int) {
	runtime.GOMAXPROCS(cores)
	n := 100000000
	A := make([]doublepair, n)
	for i := 0; i < n; i++ {
		A[i].x = rand.Float64()
		A[i].y = rand.Float64()
	}

	start := time.Now()

	sort.Sort(pairs(A))
	fmt.Println("********** Result for Builtin Sort **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Check the array is sorted: ", sort.IsSorted(pairs(A)))
}



