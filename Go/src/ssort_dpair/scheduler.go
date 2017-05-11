package ssort_dpair

import (
	"sort"
	"runtime"
	"sync"
)

type Continuation func(tuple)

func schedule_sort(A sort.Interface, i int, j int, cpus int) {
	var enqueue Continuation

	if cpus == 0 {
        cpus = runtime.GOMAXPROCS(runtime.NumCPU())
    }
    
	if cpus == 1 { Qsort_serial(A, 0, A.Len()) }

	wg := new(sync.WaitGroup)
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
			qsort_worker(A, t.x, t.y, enqueue)
			wg.Done()
		}
	}
	for td := 0; td < cpus; td++ {
		go func() {
			for t := range queue {
				qsort_worker(A, t.x, t.y, enqueue)
				wg.Done()
			}
		}()
	}

	qsort_worker(A, i, j, enqueue)
	wg.Wait()
	close(queue)
}
