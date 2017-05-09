package main

import (
	"fmt"
	"qosort"
	"os"
	"runtime"
	"strconv"
)

func main() {
	runtime.GOMAXPROCS(40)
	cores := runtime.GOMAXPROCS(runtime.NumCPU())

	args := os.Args
	n := 1000000
	var err error
	if len(args) > 1 {
		n, err = strconv.Atoi(args[1])
		if err != nil { fmt.Println("Need input argument for array length n.") }
	}

	qosort.Test_qsort_parallel(cores, n, 5)
}

