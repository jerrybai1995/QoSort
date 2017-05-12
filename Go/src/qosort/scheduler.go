package qosort

import (
	"sort"
	"runtime"
	"sync"
)

type Continuation func(tuple)  // More like a callback

func schedule_sort(A sort.Interface, i int, j int) {
	var enqueue Continuation
	cpus := runtime.GOMAXPROCS(runtime.NumCPU())
	if cpus == 1 { Qsort_serial(A, 0, A.Len()) }

	// Set up a wait group for synchronization purposes
	wg := new(sync.WaitGroup)

	// Since we have cpus # of cores, we set up a queue with size 2 times as large
	queue := make(chan tuple, cpus * 2)

	enqueue = func(t tuple) {
		if t.y - t.x < SEQ_SORT_THRESHOLD {
			qsort_worker(A, t.x, t.y, enqueue)
			return
		}
		wg.Add(1)
		select {
		case queue <- t:
		default:
			// Critical: if the queue is full (so blocked), the select will fall into this case so that the worker's
			// resources is best used!
			qsort_worker(A, t.x, t.y, enqueue)
			wg.Done()
		}
	}
	for td := 0; td < cpus; td++ {
		// Run worker pool
		go func() {
			for t := range queue {
				qsort_worker(A, t.x, t.y, enqueue)
				wg.Done()
			}
		}()
	}

	// Pushing the initial task!
	qsort_worker(A, i, j, enqueue)
	wg.Wait()
	close(queue)
}
