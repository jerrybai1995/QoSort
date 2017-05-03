## QoSort

>A parallel sorting library in Go.

### Introduction

--
Currently our algorithm implements a work-balanced parallel version of the well-known quicksort algorithm, with tweaks to improve performance with regard to concurrency management of Golang. 



### Author

- Shaojie Bai (shaojieb@andrew.cmu.edu)
- Yutong Chen (yutongc@andrew.cmu.edu)



### Usage

--

#### To Compile

To compile our library on a UNIX machine, you need to make sure that Golang is correctly installed (recommended version 1.5+). Then, untar the file `qosort.tar` and navigate to the working directory:

```sh
$ tar xvf qosort.tar
$ cd qosort
$ make
```

**Note**: The `Makefile` sets the environment variable `GOPATH` to the `[ABSOLUTE_PATH]/qosort` working directory. It also creates an executable called `run`, which already exists in the tarball. 

#### To Run

To run the parallel sort test on an array of size `length`, do

```sh
$ ./run [length]
```

For example, running `./run 100000000` will run the sorting experiment on 100M pairs of double. It shall run the sorting task for 5 times and output the best time out of the 5.

#### Note

If you are running with Go with version <1.5, it may automatically set the default max # of processors the library are able to leverage to 1. But we should have solved this problem for most of the machines and Golang versions before 1.5.



### Algorithm

--
#### Pivot selection

Our current pivot selection method samples a range of elements chosen evenly distributed in the array, then choose the median to be the pivot. The sample size is either 15 or 9 depending on input size. The sample sizes and thresholds are optimized base on our testing results.

Other tested pivot and divide-and-conquer methods include randomly selecting one element as pivod, randomly sampling three elements and choose the median. Spliting array into three sections and recruse on each section. All of these methods were outperformed by our current method in test setting. These alternative implementations are left in the code and can be tested if input or hardware setting changes.

#### Granularity

Granularity decisions include threshold to switch to serial version of quicksort (when input size doesn't justify parallel overhead anymore) and threshold to switch to insertion sort (for base cases). We expect these two threshold to be dependent on the machine hardware, and current threshold are determined base on testing on AFS UNIX machines. Thresholds are defined as globals and easily modifiable.

#### Support for generic types

Support for **generic** type is achieved using Golang's [``sort.Interface``](https://Golang.org/pkg/sort/#Interface) structure, users must implement this interface for their specific datatype for it to be compatible with our sorter.



### Parallelization

--
#### Worker pool

Worker pool is set-up by spawning ``cpus`` number of goroutines, where ``cpus`` is the number of logical cores available on the machine, given by ``runtime.NumCPU()``. After initial boot up, work queue and inter-worker communication is handled by a channel. Each call to sort method function would take in a continuation function, which is used to enqueue the next task or keep working on the task in case queue is full.

#### Reducing communication traffic

Opposing to using a dedicated work dispatcher to handle work queue and communication, we chose to use only one channel as a work queue to reduce unnecessary communication. The channel is set up upon starting the program, and access given to each worker goroutine. After that, each worker routine can fetch and push tasks to the queue channel directly. The fetch process is implemented using Golang's built-in ``range chan`` method for simplicity (so that we don't need to worry about the atomicity of the fetch operation).



### Results

--

On a 20-core machine (e.g. UNIX 5, with 40 logical processors), our implementation is able to sort 100 million pairs of double (float64) based on the key within 5.9 seconds (on average). For example:

```
$ ./run 100000000
Run 0: Starting to initialize random double-pair array...
Run 0: Starting to execute
Run 0: Finished execution. Check the array is sorted: true.

Run 1: Starting to initialize random double-pair array...
Run 1: Starting to execute
Run 1: Finished execution. Check the array is sorted: true.

Run 2: Starting to initialize random double-pair array...
Run 2: Starting to execute
Run 2: Finished execution. Check the array is sorted: true.

Run 3: Starting to initialize random double-pair array...
Run 3: Starting to execute
Run 3: Finished execution. Check the array is sorted: true.

Run 4: Starting to initialize random double-pair array...
Run 4: Starting to execute
Run 4: Finished execution. Check the array is sorted: true.

********** Result for 5 runs of (Optimized) Parallel Quicksort **********
Number of processors used:  40
Time elapsed best of 5 (seconds):  5.778199958
```

The built-in Golang's sort, [`sort.Sort`](https://golang.org/pkg/sort/#Sort), is able to sort 100M such pairs in about 69 seconds (so we achieved a ~12x speedup) with a 20-core machine (UNIX 5).

