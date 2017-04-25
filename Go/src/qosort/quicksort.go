package qosort

import (
	"sort"
	"sync"
)

var ISORT_THRESHOLD = 12

func QuickSort(A sort.Interface) {
	qsort(A, 0, A.Len())
}

func qsort(A sort.Interface, i int, j int) {
	wg := new(sync.WaitGroup)
	if (j - i) < ISORT_THRESHOLD {
		insertion_sort(A, i, j)
	} else {
		mid := split2(A, i, j)
		wg.Add(1)
		go func() {
			qsort(A, i, mid)
			wg.Done()
		}()
		qsort(A, mid, j)
	}
	wg.Wait()
}

func qsort_by3(A sort.Interface, i int, j int) {
	wg := new(sync.WaitGroup)
	if (j - i) < ISORT_THRESHOLD {
		insertion_sort(A, i, j)
	} else {
		L, M, mid_exist := split3(A, i, j)
		wg.Add(1)
		go func() {
			qsort(A, i, L)
			wg.Done()
		}()
		if mid_exist {
			wg.Add(1)
			go func() {
				qsort(A, L, M)
				wg.Done()
			}() }
		qsort(A, M, j)

	}
	wg.Wait()
}

func qsort_serial(A sort.Interface, i int, j int) {
	n := j - i
	for n > ISORT_THRESHOLD {
		L, M, mid_exist := split3(A, i, j)
		if mid_exist { qsort_serial(A, L, M) }
		qsort_serial(A, M, j)
		n = L - i
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
	sort5(A, i, j)
	A.Swap(i, i+2)
	A.Swap(i+4, j-1)  // To maintain invariant that R is larger than pivot and L is smaller than pivot
	pivot := i
	L, R := pivot+1, j-1

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

	A.Swap(i, i+3)
	A.Swap(i+1, i+6)
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