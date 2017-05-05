package qosortv2

func insertion_sort(A []qselem, lo int, hi int) {
    for i := lo; i < hi; i++ {
        k := i
        for k > lo && A[k].Less(A[k-1]) {
            A[k], A[k-1] = A[k-1], A[k]
            k--
        } 
    }
}