package qosort

import (
	"fmt"
	"sort"
	"time"
	"math/rand"
	"runtime"
	"math"
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

func Test_qsort_parallel(cores int, length int, times int) {
	runtime.GOMAXPROCS(cores)
	n := length
	exe_time := 10000.0
	for t := 0; t < 5; t++ {
		fmt.Printf("Run %d: Starting to initialize random double-pair array...\n", t)
		A := make([]doublepair, n)
		for i := 0; i < n; i++ {
			A[i].x = rand.Float64()
			A[i].y = rand.Float64()
		}
		fmt.Printf("Run %d: Starting to execute\n", t)
		start := time.Now()
		Qsort_parallel(pairs(A), 0, n)
		diff := time.Since(start).Seconds()
		fmt.Printf("Run %d: Finished execution. Check the array is sorted: %t.\n\n", t, sort.IsSorted(pairs(A)))
		exe_time = math.Min(exe_time, diff)
	}

	fmt.Println("********** Result for 5 runs of (Optimized) Parallel Quicksort **********")
	fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
	fmt.Println("Time elapsed best of 5 (seconds): ", exe_time)
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



