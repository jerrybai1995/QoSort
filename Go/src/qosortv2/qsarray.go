package qosortv2

type qsarray []qselem

func (A qsarray) Len() int {
    return len(A)
}

func (A qsarray) Swap(i,j int) {
    A[i], A[j] = A[j], A[i]
}

func (A qsarray) Less(i,j int) bool {
    return A[i].Less(A[j])
}