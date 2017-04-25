package qosort

import (
	"sync"
	"sort"
)

var ISORT_THRESHOLD = 20

func QuickSort(A sort.Interface, i int, j int) {
	qsort(A, i, j)
}

func qsort(A sort.Interface, i int, j int) {
	var wg sync.WaitGroup

	if (j - i) < ISORT_THRESHOLD {
		insertion_sort(A, i, j)
	} else {
		L, M, mid_exist := split3(A, i, j)
		wg.Add(1)
		go func(){
			qsort(A, i, L)
			wg.Done()
		}()
		if mid_exist {
			wg.Add(1)
			go func(){
				qsort(A, L, M)
				wg.Done()
			}()
		}
		qsort(A, M, j)
		wg.Wait()

	}
}

func qsort_serial(A sort.Interface, i int, j int) {
	n := j - i
	for n > 24 {
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


func split3(A sort.Interface, i int, j int) (int, int, bool) {
	sort5(A, i, j)

	p1, p2 := i+1, i+3
	if !A.Less(i, i+1) { p1 = p2 }
	if !A.Less(i+3, i+4) { p2 = p1 }

	L, R := i, j-1
	for A.Less(L, p1) { L++ }
	for A.Less(p2, R) { R-- }
	M := L
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
	return L, M, A.Less(p1, p2)
}