package ssort_dpair

import (
	"fmt"
	"sort"
	"time"
	"math/rand"
	"runtime"
	"math"
)

func Test_ssort_parallel(cores int, length int, times int) {
    runtime.GOMAXPROCS(cores)
    n := length
    exe_time := 10000.0
    for t := 0; t < times; t++ {
        fmt.Printf("\n\nRun %d: Starting to initialize random double-pair array of size %d...\n", t, length)
        A := make_random_doublepairs(n)

        fmt.Printf("Run %d: Starting to execute\n", t)
        start := time.Now()
        SampleSort(A)
        diff := time.Since(start).Seconds()
        fmt.Printf("Run %d: Finished execution. Check the array is sorted: %t.\n\n", t, sort.IsSorted(pairs(A)))
        exe_time = math.Min(exe_time, diff)
    }

    fmt.Println("********** Result for 5 runs of (Optimized) Parallel Samplesort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed best of 5 (seconds): ", exe_time)
}

func Test_qsort_parallel(cores int, length int, times int) {
    runtime.GOMAXPROCS(cores)
    n := length
    exe_time := 10000.0
    for t := 0; t < times; t++ {
        fmt.Printf("\n\nRun %d: Starting to initialize random double-pair array of size %d...\n", t, length)
        A := make_random_doublepairs(n)

        fmt.Printf("Run %d: Starting to execute\n", t)
        start := time.Now()
        Qsort_parallel(pairs(A), 0, n, 0)
        diff := time.Since(start).Seconds()
        fmt.Printf("Run %d: Finished execution. Check the array is sorted: %t.\n\n", t, sort.IsSorted(pairs(A)))
        exe_time = math.Min(exe_time, diff)
    }

    fmt.Println("********** Result for 5 runs of (Optimized) Parallel Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed best of 5 (seconds): ", exe_time)
}


func Test_Qsort_serial(length int) {
    fmt.Printf("Testing Qsort_serial sort with input size: %d\n", length)

    n := length

    A := make_random_doublepairs(n)

    start := time.Now()
    Qsort_serial(pairs(A), 0, n)

    fmt.Println("********** Result for Serial Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed: ", time.Since(start))
    fmt.Printf("Finished execution. Check the array is sorted: %t.\n\n", sort.IsSorted(pairs(A)))
}

func Test_Qsort_naive_par(cores int, length int) {
    runtime.GOMAXPROCS(cores)

    fmt.Printf("Testing Qsort_serial sort with input size: %d\n", length)

    n := length

    A := make_random_doublepairs(n)

    start := time.Now()
    Qsort_naive_parallel(pairs(A), 0, n)
    
    fmt.Println("********** Result for naive-par Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed: ", time.Since(start))
    fmt.Printf("Finished execution. Check the array is sorted: %t.\n\n", sort.IsSorted(pairs(A)))    

}


func make_random_doublepairs(n int) []doublepair {
    A := make([]doublepair, n)
    for i := 0; i < n; i++ {
        e := doublepair{rand.Float64(), rand.Float64()}
        A[i] = e
    }
    return pairs(A)
}