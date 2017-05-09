package main

import (
    "qosortv2"
)


func main() {

    // cores := runtime.GOMAXPROCS(runtime.NumCPU())
    // n := 100000000
    // times := 5

    // qosortv2.Test_insertion_sort(100000)
    // qosortv2.Test_Qsort_serial(100000000)
    // qosortv2.Test_Qsort_naive_par(4,10000000)

    qosortv2.Test_ssort_parallel(8, 1000000, 1)
    qosortv2.Test_qsort_parallel(8, 1000000, 1)

}
