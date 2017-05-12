package main

import (
	"ssort_dpair"
)

func main() {

    // cores := runtime.GOMAXPROCS(runtime.NumCPU())
    // n := 100000000
    // times := 5

    // ssort_dpair.Test_Qsort_serial(10000000)
    // ssort_dpair.Test_Qsort_naive_par(4,10000000)

    ssort_dpair.Test_ssort_parallel(40, 100000000, 1)
    ssort_dpair.Test_qsort_parallel(40, 100000000, 1)

}





