package main

import (
	"fmt"
	"qosort"
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


func main() {

	//qosort.Test_qsort_serial(10)
	//qosort.Test_qsort_by3(10)
	//qosort.Test_qsort_parallel(10)
	//qosort.Test_sort(10)
	//
	//qosort.Test_qsort_serial(20)
	//qosort.Test_qsort_by3(20)
	//qosort.Test_qsort_parallel(20)
	//qosort.Test_sort(20)
	//
	//qosort.Test_qsort_serial(30)
	//qosort.Test_qsort_by3(30)
	//qosort.Test_qsort_parallel(30)
	//qosort.Test_sort(30)
    //
	//qosort.Test_qsort_serial(40)
	//qosort.Test_qsort_by3(40)
	qosort.Test_qsort_parallel(40)
	qosort.Test_qsort_parallel_ref(40)
	qosort.Test_sort(40)
}


