




###Introduction
--
Currently our algorithm implements a parallel version of the well-known quicksort algorithm, with tweaks to improve performance with regard to concurrecy management of GoLang. 

###Usage
--
TODO

###Algorithm
--
#### Pivot selection

Our current pivot selection method samples a range of elements chosen evenly distributed in the array, then choose the median to be the pivot. The sample size is either 15 or 9 depending on input size. The sample sizes and thresholds are optimized base on our testing results.

Other tested pivot and divide-and-conquer methods include randomly selecting one element as pivod, randomly sampling three elements and choose the median. Spliting array into three sections and recruse on each section. All of these methods were outperformed by our current method in test setting. These alternative implementations are left in the code and can be tested if input or hardware setting changes.

#### Granularity
Granularity decisions include threshold to switch to serial version of quicksort (when input size doesn't justify parallel overhead anymore) and threshold to switch to insertion sort (for base cases). We expect these two threshold to be dependent on the machine hardware, and current threshold are determined base on testing on AFS unix machines. Thresholds are defined as globals and easily modifiable.


#### Support for generic types
Support for generic type is achieved using GoLang's [``sort.Interface``](https://golang.org/pkg/sort/#Interface) structure, users must implement this interface for their specific datatype for it to be comptible with our sorter.



###Parallelization
--
#### Worker pool
Worker pool is set-up by spawning ``cpus`` number of goroutines, where ``cpus`` is the number of logical cores available on the machine, given by ``runtime.NumCPU()``. After initial boot up, work queue and inter-worker communication is handled by a channel. Each call to sort method function would take in a continuition function, which is used to enqueue the next task or keep working on the task in case queue is full.


#### Reducing communication traffic
Opposing to using a dedicated work dispatcher to handle work queue and communication, we chose to use only one channel as a work queue to reduce unnecessary communication. The channel is set up upon starting the program, and access given to each worker goroutine. After that, each worker routine can fetch and push tasks to the queue channel directly. The fetch process is implemented using GoLang's built-in ``range chan`` method for simplicity.
