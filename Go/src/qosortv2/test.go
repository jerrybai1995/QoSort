package qosortv2

import (
    "fmt"
    "math"
    "math/rand"
    "runtime"
    "sort"
    "time"
)

func Test_ssort_parallel(cores int, length int, times int) {
    runtime.GOMAXPROCS(cores)
    n := length
    exe_time := 10000.0
    for t := 0; t < times; t++ {
        fmt.Printf("Run %d: Starting to initialize random double-pair array of size %d...\n", t, length)
        A := make_random_doublepairs(n)

        fmt.Printf("Run %d: Starting to execute\n", t)
        start := time.Now()
        SampleSort(A)
        diff := time.Since(start).Seconds()
        fmt.Printf("Run %d: Finished execution. Check the array is sorted: %t.\n\n", t, sort.IsSorted(qsarray(A)))
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
        fmt.Printf("Run %d: Starting to initialize random double-pair array...\n", t)
        A := make_random_doublepairs(n)

        fmt.Printf("Run %d: Starting to execute\n", t)
        start := time.Now()
        Qsort_parallel(A, 0, n)
        diff := time.Since(start).Seconds()
        fmt.Printf("Run %d: Finished execution. Check the array is sorted: %t.\n\n", t, sort.IsSorted(qsarray(A)))
        exe_time = math.Min(exe_time, diff)
    }

    fmt.Println("********** Result for 5 runs of (Optimized) Parallel Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed best of 5 (seconds): ", exe_time)
}

func Test_insertion_sort(length int) {
    fmt.Printf("qosortv2: Testing insertion_sort with input size: %d \n", length)

    n := length

    A := make_random_doublepairs(n)

    start := time.Now()
    insertion_sort(A, 0, n)
    diff := time.Since(start).Seconds()
    fmt.Printf("Finished execution. Check the array is sorted: %t.\n\n", sort.IsSorted(qsarray(A)))
    fmt.Println("Time elapsed (seconds): ", diff)
}

func Test_Qsort_serial(length int) {
    fmt.Printf("Testing Qsort_serial sort with input size: %d\n", length)

    n := length

    A := make_random_doublepairs(n)

    start := time.Now()
    Qsort_serial(A, 0, n)

    fmt.Println("********** Result for Serial Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed: ", time.Since(start))
    fmt.Printf("Finished execution. Check the array is sorted: %t.\n\n", sort.IsSorted(qsarray(A)))
}

func Test_Qsort_naive_par(cores int, length int) {
    runtime.GOMAXPROCS(cores)

    fmt.Printf("Testing Qsort_serial sort with input size: %d\n", length)

    n := length

    A := make_random_doublepairs(n)

    start := time.Now()
    Qsort_naive_parallel(A, 0, n)
    
    fmt.Println("********** Result for naive-par Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed: ", time.Since(start))
    fmt.Printf("Finished execution. Check the array is sorted: %t.\n\n", sort.IsSorted(qsarray(A)))    

}

func make_random_doublepairs(n int) []qselem {
    A := make([]qselem, n)
    for i := 0; i < n; i++ {
        e := doublepair{rand.Float64(), rand.Float64()}
        A[i] = e
    }
    return A
}
