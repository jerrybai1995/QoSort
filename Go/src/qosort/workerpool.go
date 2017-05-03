package qosort
// package main

import (
        "fmt"
        "math/rand"
        "runtime"
        "sort"
        "sync"
        "time"
)

type CallBack func()
type SortMethod func(sort.Interface, int, int, chan Job, CallBack)

type Job struct {
    lo int
    hi int
    data sort.Interface
}

type Worker struct {
    sortMethod SortMethod
    workerPool chan *Worker
    jobQueue   chan Job
    nextJob    chan Job
    stop       chan bool
    wg         *sync.WaitGroup
}

func (w *Worker) start() {
    go func() {
        var j Job
        for {

            w.workerPool <- w

            select {
            case j = <-w.nextJob:
                w.sortMethod(j.data, j.lo, j.hi, w.jobQueue,
                    func(){
                        w.wg.Add(1)
                        })
                w.wg.Done()

            case stop := <-w.stop:
                if stop {
                    w.stop <- true
                    return
                }
            }
        }
    }()
}


func newWorker(pool chan *Worker, jQueue chan Job, sMethod SortMethod, wg *sync.WaitGroup) *Worker {
    return &Worker{
        sortMethod: sMethod,
        workerPool: pool,
        jobQueue:   jQueue,
        nextJob:    make(chan Job),
        stop:       make(chan bool),
        wg:         wg,
    }
}


type Dispatcher struct {
    workerPool chan *Worker
    jobQueue   chan Job
    stop       chan bool
    wg         *sync.WaitGroup
}

func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-d.jobQueue:
            worker := <-d.workerPool
            worker.nextJob <- job

        case stop := <-d.stop:
            if stop {
                for i := 0; i < cap(d.workerPool); i++ {
                    worker := <-d.workerPool
                    worker.stop <- true
                    <-worker.stop
                }
                d.stop <- true
                return
            }
        }
    }
}

func newDispatcher(workerPool chan *Worker, jQueue chan Job, sMethod SortMethod, wg *sync.WaitGroup) *Dispatcher {
    d := &Dispatcher{
        workerPool: workerPool,
        jobQueue:   jQueue,
        stop:       make(chan bool),
        wg:         wg,
    }

    for i := 0; i < cap(d.workerPool); i++ {
        worker := newWorker(d.workerPool, d.jobQueue, sMethod, wg)
        worker.start()
    }

    go d.dispatch()
    return d
}

type Pool struct {
    jobQueue    chan Job
    dispatcher  *Dispatcher
    wg          *sync.WaitGroup
    sortMethod  SortMethod
}

func NewPool(numWorkers int, jobQueueLen int, sMethod SortMethod) *Pool {
    jobQueue := make(chan Job, jobQueueLen)
    workerPool := make(chan *Worker, numWorkers)
    wg := new(sync.WaitGroup)

    pool := &Pool{
        jobQueue:   jobQueue,
        dispatcher: newDispatcher(workerPool, jobQueue, sMethod, wg),
        sortMethod: sMethod,
        wg: wg,
    }

    return pool
}

func (p *Pool) JobDone() {
    p.wg.Done()
}

func (p *Pool) WaitCount(count int){
    p.wg.Add(count)
}

func (p *Pool) WaitAll() {
    p.wg.Wait()
}

func (p *Pool) ShutDown() {
    p.dispatcher.stop <- true
    <-p.dispatcher.stop
}


/**************************************************
    Sample SortMethod impl
*/
 

var ISORT_THRESHOLD = 12

// Regular quicksort with carefully picked median, and executed recursively
func qsort_serial(A sort.Interface, i int, j int, jQ chan Job, f CallBack) {
    n := j - i
    for n > ISORT_THRESHOLD {
        mid := split2(A, i, i + n)
        j := Job{
            lo:   mid,
            hi:   i+n,
            data: A,
        }
        select {
        case jQ <- j:
            f()
        default:
            qsort_serial(A, mid, i+n, jQ, f)
        }
        n = mid - i
    }
    insertion_sort(A, i, i + n)
}

func insertion_sort(A sort.Interface, lo int, hi int) {
    for i := lo; i < hi; i++ {
        k := i
        for k > lo && A.Less(k, k-1) {
            A.Swap(k, k-1)
            k--
        }
    }
}

func sort5(A sort.Interface, i int, j int) {
    size := 5
    m := (j - i) / (size + 1)
    for l := 0; l < size; l++ {
        A.Swap(i+l, i+m*(l+1))
    }
    insertion_sort(A, i, i + size)
}

func split2(A sort.Interface, i int, j int) int {
    sort5(A, i, j)
    A.Swap(i, i+2)
    A.Swap(i+4, j-1)  // To maintain invariant that R is larger than pivot and L is smaller than pivot
    pivot := i
    L, R := pivot+1, j-1

    for A.Less(L+1, pivot) { L++ }
    for A.Less(pivot, R-1) { R-- }

    M := L
    for {
        for M < R && !A.Less(pivot, M+1) { M++ }  // if A[M+1] <= A[pivot], M++
        for M < R && A.Less(pivot, R-1) { R-- }   // if A[pivot] < A[R-1], R--
        if R - M <= 1 {
            // A[M] <= pivot and yet A[R] > pivot
            break
        }
        A.Swap(M+1, R-1)
        M++
        R--
    }
    A.Swap(i, M)

    return R
}

func split3(A sort.Interface, i int, j int) (int, int, bool) {
    sort5(A, i, j)

    p1, p2 := i, i+1

    A.Swap(i, i+1)
    A.Swap(i+1, i+3)
    mid_exist := A.Less(p1, p2)

    L, R := i+2, j-1
    for A.Less(L, p1) { L++ }
    for A.Less(p2, R) { R-- }
    M := L

    // A[i <= x < L] < pivot1   i...L-1
    // A[R < x <= j] > pivot2   L...M-1
    for M <= R {
        if A.Less(M, p1) {
            // Should have been in the first 1/3
            A.Swap(M, L)
            L++
        } else if A.Less(p2, M) {
            if A.Less(R, p1) {
                A.Swap(M, L)
                A.Swap(R, L)
                L++
            } else {
                A.Swap(M, R)
            }
            R--
            for A.Less(p2, R) { R-- }
        }
        M++
    }
    A.Swap(i, L-2)
    A.Swap(i+1, L-1)
    A.Swap(L-1, M-1)
    L = L-2
    M = M-1

    return L, M, mid_exist
}

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

/***************************************************
    Sample usage
*/

func main() {
    cores := runtime.GOMAXPROCS(0)

    numWorkers := cores
    pool := NewPool(numWorkers, 100, qsort_serial)

    n := 10000000
    A := make([]doublepair, n)
    for i := 0; i < n; i++ {
        A[i].x = rand.Float64()
        A[i].y = rand.Float64()
    }

    j := Job{
        lo: 0,
        hi: n,
        data: pairs(A),
    }

    start := time.Now()

    pool.WaitCount(1)
    pool.jobQueue <- j
    pool.WaitAll()
    pool.ShutDown()

    fmt.Println("********** Result for Serial Quicksort **********")
    fmt.Println("Number of processors used: ", runtime.GOMAXPROCS(0))
    fmt.Println("Time elapsed: ", time.Since(start))
    fmt.Println("Check they array is sorted: ", sort.IsSorted(pairs(A)))

}



