package main

import (
    // "fmt"
    "qosortv2"
    // "runtime"
    // "strconv"
)

func main() {

    // cores := runtime.GOMAXPROCS(runtime.NumCPU())
    // n := 100000000
    // times := 5

    // qosortv2.Test_qsort_parallel(cores, n, times)

    // qosortv2.Test_insertion_sort(100000)
    // qosortv2.Test_Qsort_serial(100000000)
    qosortv2.Test_Qsort_naive_par(4,10000000)

}
