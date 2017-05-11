package ssort_dpair

import (
	"sort"
	"sync"
)

var SEQ_SORT_THRESHOLD = 120
var ISORT_THRESHOLD = 12

func QuickSort(A sort.Interface) {
	Qsort_naive_parallel(A, 0, A.Len())
}

// median_of_three moves the median of the three to the first memory location m1
func median_of_three(data sort.Interface, m1, m0, m2 int) {
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	// data[m0] <= data[m1]
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}

// The regular, serial quicksort to be scheduled by a worker (runner)
func qsort_worker(A sort.Interface, i int, j int, f Continuation) {
	n := j - i
	for n > SEQ_SORT_THRESHOLD {
		mid := split2(A, i, i + n)
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
		Qsort_serial(A, i, i + n)
	} else {
		insertion_sort(A, i, i + n)
	}
}

// Parallel quicksort with optimized workload
func Qsort_parallel(A sort.Interface, i int, j int, cores int) {
	schedule_sort(A, i, j, cores)
}

// Regular quicksort with carefully picked median, and parallelized by goroutines at each recursive call
func Qsort_naive_parallel(A sort.Interface, i int, j int) {
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

// The quicksort that, instead of splitting the original array into two parts, divides it into 3 sections
func qsort_by3(A sort.Interface, i int, j int) {
	if (j - i) < SEQ_SORT_THRESHOLD {
		insertion_sort(A, i, j)
	} else {
		L, M, mid_exist := split3(A, i, j)
		Qsort_naive_parallel(A, i, L)
		if mid_exist { Qsort_naive_parallel(A, L, M) }
		Qsort_naive_parallel(A, M, j)

	}
}

// Regular quicksort with carefully picked median, and executed recursively
func Qsort_serial(A sort.Interface, i int, j int) {
	n := j - i
	for n > ISORT_THRESHOLD {
		mid := split2(A, i, i + n)
		Qsort_serial(A, mid, i + n)
		n = mid - i
	}
	insertion_sort(A, i, i + n)
}

func insertion_sort(A sort.Interface, lo int, hi int) {
	for i := lo; i < hi; i++ {
		k := i
		for k > lo && A.Less(k, k-1) {
			A.Swap(k, k-1)
			k--
		}
	}
}

func sort5(A sort.Interface, i int, j int) {
	size := 5
	m := (j - i) / (size + 1)
	for l := 0; l < size; l++ {
		A.Swap(i+l, i+m*(l+1))
	}
	insertion_sort(A, i, i + size)
}

func split2(A sort.Interface, i int, j int) int {
	// sort5(A, i, j)
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
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (j - i) / 8
		median_of_three(A, i, i+s, i+2*s)
		median_of_three(A, m, m-s, m+s)
		median_of_three(A, j-1, j-1-s, j-1-2*s)
	}
	median_of_three(A, i, m, j-1)

	// To maintain invariant that R is larger than pivot and L is smaller than pivot
	pivot := i
	L, R := pivot, j-1

	for A.Less(L+1, pivot) { L++ }
	for A.Less(pivot, R-1) { R-- }

	M := L
	for {
		for M < R && !A.Less(pivot, M+1) { M++ }  // if A[M+1] <= A[pivot], M++
		for M < R && A.Less(pivot, R-1) { R-- }   // if A[pivot] < A[R-1], R--
		if R - M <= 1 {
			// A[M] <= pivot and yet A[R] > pivot
			break
		}
		A.Swap(M+1, R-1)
		M++
		R--
	}
	A.Swap(i, M)

	return R
}

func split3(A sort.Interface, i int, j int) (int, int, bool) {
	sort5(A, i, j)

	p1, p2 := i, i+1

	A.Swap(i, i+1)
	A.Swap(i+1, i+3)
	mid_exist := A.Less(p1, p2)

	L, R := i+2, j-1
	for A.Less(L, p1) { L++ }
	for A.Less(p2, R) { R-- }
	M := L

	// A[i <= x < L] < pivot1   i...L-1
	// A[R < x <= j] > pivot2   L...M-1
	for M <= R {
		if A.Less(M, p1) {
			// Should have been in the first 1/3
			A.Swap(M, L)
			L++
		} else if A.Less(p2, M) {
			if A.Less(R, p1) {
				A.Swap(M, L)
				A.Swap(R, L)
				L++
			} else {
				A.Swap(M, R)
			}
			R--
			for A.Less(p2, R) { R-- }
		}
		M++
	}
	A.Swap(i, L-2)
	A.Swap(i+1, L-1)
	A.Swap(L-1, M-1)
	L = L-2
	M = M-1

	return L, M, mid_exist
}