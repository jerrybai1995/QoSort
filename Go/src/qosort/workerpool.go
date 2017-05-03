package qosort

import (
        "sort"
        "sync"
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

func qsort_runner(A sort.Interface, i int, j int, jQ chan Job, f CallBack) {
    n := j - i
    for n > SEQ_SORT_THRESHOLD {
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
            qsort_runner(A, mid, i+n, jQ, f)
        }
        n = mid - i
    }
    insertion_sort(A, i, i + n)
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




