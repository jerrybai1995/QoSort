package main

import (
    // "fmt"
    "ssort_dpair"
    "runtime"
)

func main() {

    cores := runtime.GOMAXPROCS(runtime.NumCPU())
    n := 100000000
    // n := 1000000
    times := 5

    ssort_dpair.Test_qsort_parallel(cores, n, times)
    ssort_dpair.Test_ssort_parallel(cores, n, times)

}
