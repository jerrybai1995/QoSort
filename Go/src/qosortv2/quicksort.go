package qosortv2

import (
    "sync"
    // "time"
    // "fmt"
)

var SEQ_SORT_THRESHOLD = 130
var ISORT_THRESHOLD = 12


func QuickSort(A []qselem) {
    Qsort_naive_parallel(A, 0, len(A))
}

func Qsort_parallel(A []qselem, i,j int) {
    schedule_sort(A, i, j)
}

/**************************************************
 * Sorting methods
 */

// Regular quicksort with carefully picked median, and executed recursively
func Qsort_serial(A []qselem, i int, j int) {
    n := j-i
    for n > ISORT_THRESHOLD {
        mid := split2(A, i, i+n)
        Qsort_serial(A, mid, i+n)
        n = mid - i
    }
    insertion_sort(A, i, i+n)
}

// Regular quicksort with carefully picked median, and parallelized by goroutines at each recursive call
func Qsort_naive_parallel(A []qselem, i int, j int) {
    wg := new(sync.WaitGroup)
    if (j - i) < SEQ_SORT_THRESHOLD {
        insertion_sort(A, i, j)
    } else {
        mid := split2(A, i, j)
        wg.Add(1)
        go func() {
            Qsort_naive_parallel(A, i, mid)
            wg.Done()
        }()
        Qsort_naive_parallel(A, mid, j)
    }
    wg.Wait()
}

// The regular, serial quicksort to be scheduled by a worker (runner)
func qsort_worker(A []qselem, i int, j int, f Continuation) {
    n := j-i
    for n > SEQ_SORT_THRESHOLD {
        mid := split2(A, i, i+n)
        if mid - i < i + n - mid {
            f(tuple{i, mid})
            i = mid
        } else {
            f(tuple{mid, i + n})
            j = mid
        }
        n = j - i
    }
    if n > 7 {
        Qsort_serial(A, i, i+n)
    } else {
        insertion_sort(A, i, i+n)
    }
}

/**************************************************
 * Sorting helper functions
 */

func median_of_three(A []qselem, m1, m0, m2 int) {
    if A[m1].Less(A[m0]) {
        swap(A, m1, m0)
    }

    if A[m2].Less(A[m1]) {
        swap(A, m2, m1)

        if A[m1].Less(A[m0]) {
            swap(A, m1, m0)
        }
    }
}

func split2(A []qselem, i, j int) int {

    m := i + (j - i)/2 // Written like this to avoid integer overflow.
    if j - i > 80 {
        s := (j - i) / 8
        r := (j - i) / 16
        median_of_three(A, i, i+s, i+2*s)
        median_of_three(A, i, i+r, i+s+r)
        median_of_three(A, m, m-s, m+s)
        median_of_three(A, m, m-s+r, m+s+r)
        median_of_three(A, j-1, j-1-s, j-1-2*s)
        median_of_three(A, j-1, j-1-r, j-1-s-r)
    } else if j - i > 40 {
        s := (j - i) / 8
        median_of_three(A, i, i+s, i+2*s)
        median_of_three(A, m, m-s, m+s)
        median_of_three(A, j-1, j-1-s, j-1-2*s)
    }
    median_of_three(A, i, m, j-1)

    pivot := A[i]
    L, R := i, j-1

    for A[L+1].Less(pivot) {L++}
    for pivot.Less(A[R-1]) {R--} 

    M := L
    for {
        for M < R && !pivot.Less(A[M+1]) {M++}
        for M < R && pivot.Less(A[R-1]) {R--}

        if R - M <= 1 {
            // A[M] <= pivot and yet A[R] > pivot
            break
        }
        swap(A, M+1, R-1)
        M++
        R--
    }
    swap(A, i, M)

    return R

}

/**************************************************
 * Array utility functions
 */

func swap(A []qselem, i,j int){
    A[i], A[j] = A[j], A[i]
}


