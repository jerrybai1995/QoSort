package qosort

import (
	"sort"
	"runtime"
	"sync"
)

func schedule_sort(A sort.Interface) {
	var enqueue func(t tuple)
	cpus := runtime.GOMAXPROCS(runtime.NumCPU())
	if cpus == 1 { qsort_serial(A, 0, A.Len()) }
	wg := new(sync.WaitGroup)
	queue := make(chan tuple, cpus * 2)
	enqueue = func(t tuple) {
		if t.y - t.x < ISORT_THRESHOLD { insertion_sort(A, t.x, t.y) }
		wg.Add(1)
		select {
		case queue <- t:
		default:
			qsort_qsub(A, t.x, t.y, enqueue)
			wg.Done()
		}
	}
	for td := 0; td < cpus; td++ {
		go func() {
			for t := range queue {
				qsort_qsub(A, t.x, t.y, enqueue)
				wg.Done()
			}
		}()
	}
}
