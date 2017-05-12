---
layout: post
title: Final Report
description: A complete report on QoSort
image: assets/images/engineering.jpg
---

## I. Summary

**QoSort** is a *parallel*, *generic* sorting library in Go. Importantly, our sorting methods are entirely consistent with the `sort.Interface` that Go's 'sort' library uses. Two methods that we've been focusing on optimizing are **quicksort** and **samplesort**. Through optimizations in terms of cache locality, thread scheduling and load balance, we are able to achieve encouraging results compared to not only Go's builtin sort, but also the best generic-type sorting method that we found online.

**UPDATE**: On May 12, we presented this project as one of the **finalists** in 15-418's parallel competition! Our slides can be found [here]({{ site.url }}/QoSort/assets/pdfs/qosort_slides.pdf). 



## II. Background

Sorting is a subject that is both easy and difficult. For a generic sort, the objective is easy: given an input data structure whose elements (of arbitrary type) support ordering, we want the eventual output to be in an increasing order (as in the defined comparisons).

#### Go & Generic Data Type

Go (Golang) is a garbage-collected language that supports lightweight multi-threading. Instead of managing instruction streams as threads, Go has goroutines that has similar abstraction as fork. Moreover, Go has relatively memory allocation compared to other peer languages such as Javascript or Haskell, so avoiding taking up too much memory is critical to implementations in Go.

One important thing about sorting in Go is that Go's standard sorting library, [`sort`](https://golang.org/pkg/sort/), takes inputs of type `sort.Interface`, which exposes an abstraction as follows:

```go
type Interface interface {
        // Len is the number of elements in the collection.
        Len() int
        // Less reports whether the element with
        // index i should sort before the element with index j.
        Less(i, j int) bool
        // Swap swaps the elements with indexes i and j.
        Swap(i, j int)
}
```

Operating on the very same interface allows our library to be consistent with Go's built-ins, and thus making it a natural extension of Go. However, this also implies some serious limitations: we can assume on neither the data structure nor the data type of the input, and there is no support for the get or set method at an index \\( i \\).

Finally, Go does not support explicit hardware scheduling (high level of abstraction, unlike pthreads) or parallel APIs such as OpenMP or ISPC. This makes efficient utilization of the cores especially important!



#### Quicksort

The idea of quicksort is two-fold:

1. Partitioning the input into two sub-sequences based on chosen pivots
2. Recursively sort on the resulting partitions

Quicksort, on average, has a time complexity of \\( O(n \log n) \\). Note that quicksort is suitable for parallel improvements because the operations take place in separate segments, and so the sorting within each segment is independent from other segments.



#### Critical Factors

A few questions that are at the core of this project:

- How to load balance?
- How to reduce memory footprint?
- How to leverage more computing resources efficiently (as well as locality)?
- How to identify the parallelism (some of which can be mutually exclusive)?

Moreover, for sorting, a general \\( O(n \log n) \\) work is inevitable. However, exploiting the multi-core machines that we have access to, we believe a sorting can be highly parallelizable if we can:

1. Identify independent work (so that the cores can run without affecting each other);
2. Identify good locality (this is especially important in sorting, since most of the time we work on data structures such as arrays);
3. Identify reasonable scheduling (i.e. we must be a lot smarter than an embarrassingly parallel solution).

These thoughts have guided us through the various optimizations and attempts that we've made throughout this project. See below (the *Approach* section) for more details!



## III. Approaches

There aren't many parallel sorting implementations (at least for the open source which we can have access to) available in developer communities such as GitHub, therefore, the specific implementations and ideas for the different approaches most came from ourselves. However, we do note that this project was initially inspired by the CMU's 15-210 sorting competition. But we aim for not merely a fast-algorithm that works on, say integers; but also a generic library that can sorts strings, floats, pairs, and basically any object that is comparable.

As is to be mentioned in the "Results" section, our method is to be compared with the best open-source Go's parallel sorting. The primary testing of our implementations was carried out on the AFS UNIX machines (Intel Xeon E5-2680 v2), which have 20 cores and are 2-way hyperthreaded. 

### QoSort.Quicksort

One of the methods that we've been focusing on optimizing is the quicksort. Two key factors at stake here:

1. Pivot selection. Especially on large inputs, this can be very important since a bad pivot not only means we would waste a lot time in each iteration doing useless linear scans, but also potential load imbalance in the subsequent passes.
2. Task scheduling. How to best organize the recursive sorting so that we can achieve load balance, good locality, and at the same time minimize unnecessary synchronization costs?

##### Quicksort v1.0

A good way to start is a naive parallelization of the quicksort task--- i.e. at each recursive call on sub-sequences, we spawn a new thread (in fact, goroutine) to handle this subarray (see **figure 1**). This will eventually turn into a recursive forking, which is subject to a number of drawbacks:

<div class="imgcap">
<img src="/QoSort/assets/images/naive.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 1: Quicksort v1.0 that is based on straight-parallelization of naive quicksort.</div>
</div>

<br />

1. Recursive call stack problem (if the input is sufficient large--- e.g. a long array of large-size elements, the stack space will be quickly consumed);
2. Load imbalance. This is something that we observed when we try to profile the execution performance. Some threads are running on large inputs while some others are not. This is due to bad pivoting as well as direct, reckless recursive fork.
3. Potential OS scheduling overhead. With the binary partitioning, the total number of goroutines will grow exponentially so that it will eventually overflow the # of cores available (provided that the input size is large enought).
4. Synchronization cost. Each parent goroutine will be responsible for synchronizing its immediate children. This "hierarchy" synchronization, we found, induces lots of cost.

**Quicksort v2.0**

Naturally, a smarter formulation of the problem is needed. We start by a smart pivoting with **median-of-three partitioning** [1]. More specifically, the idea is straightforward, the quality of the median estimate picked can be greatly improved with more sample points. 

- For inputs of length > 1600, we pick 15 uniformly spaced elements randomly from the input (i.e. at index i=0, n/15, 2n/15, ...). This will compose a series of **median-of-three** swaps using the generic `sort.Interface`'s swap and comparison operations (without explicit get and sets).
- For inputs of length in [40, 1600], we pick 9 uniformly spaced elements randomly from the input and swap the median to the front of the input.

Overall, we found this pivoting to be fast (very few computation required), and the pivot picked is usually within 4% of the median's quantile (50%). With the data of scale of ~100M, this is very impressive result. 

Moreover, we need to optimize the scheduling. One part we spotted that could lead to potential performance improvement is that we can think of a sorting as a "task". Each iteration, therefore, will be the process of breaking a task into two portions. This allows for better organization of the tasks through a task queue. Meanwhile, with the two partitions resulted, we iterate on the larger portion instead of pushing both tasks onto the queue. This permits us to better exploit cache locality, since the subsequent pass will be working on the same memory range. 

This leads to the "master-dispatch" model for sorting (see **figure 2**). One master, more specifically, is delegated the responsibility of collecting new tasks generated. Then, it will find the next idle worker and dispatch/push the task to it. (**Note**: the task does not need to be a subsequence. Instead, it will simply be a tuple of indices, and the in-place swaps by each work are independent.)

<div class="imgcap">
<img src="/QoSort/assets/images/dispatch.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 2: Quicksort v2.0 with better pivoting and master scheduling.</div>
</div>

<br />

With this improvement, we are able to achieve a good deal of speed up so that sorting 100M elements can be completed within 15-16 seconds. Nevertheless, we found some problem with this model as well:

1. Each worker must have its separate channel of communication with the master. The overhead in this is a major problem on a multi-core machine since channeling can introduce some overhead. Note that this is an overhead that quicksort v1.0 does not have (and partially accounts for v2.0's relatively suboptimal performance).
2. Synchronization. The master will have to find the next idle worker, and this requires synchronizations.

We thus proceed to the next level of optimization. 

##### Quicksort v3.0

Instead of having a central dispatcher, we can eliminate this master and instead have each worker actively pull from the task queue. This will reduce all the communication overheads that previously exists and better utilize workers' free resources. In its iteration, the worker shall keep iterate on the larger portion and push the smaller half to the task queue. However, there are a few more small optimizations:

- If the queue is full, instead of blocking, the worker that pushes will be responsible for handling this task. This is the best choice because it avoids the blocking due to slow consumer and meanwhile leverages good spatial locality since it is handled by the very same worker.
- The base case uses serial quicksort and insertion sort--- we find them quite beneficial for the performance since it is not necessary to do the scheduling on smaller inputs.

<div class="imgcap">
<img src="/QoSort/assets/images/onequeue.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 3: Quicksort v3.0. A pool of workers directly pull from a task queue. Much better load balance and communication scheme, good cache locality. Eventual result on sorting 100M double-pair elements was very encouraging.</div>
</div>

<br />

This method (see **figure 3**) has achieved very good result, being able to sort 100M generic data type (in our case, float64-pair elements) within about 4.5 seconds. But we seek to further improve the model! In general, in our measuring, this modeling of having "one task queue" can provide good load balance as well.

**Quicksort v4.0**

One observation we made is that, in the first few passes, many workers stay idle, whereas some other workers have heavy loads. For example, in the 2nd iteration (see **figure 4**), only two workers will be working (each with size n/2), while the rest of the cores are free. 

<div class="imgcap">
<img src="/QoSort/assets/images/split-in-2.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 4: Splitting-in-2 can lead to many idle workers in the initial/bootstrap phase, which wastes resources.</div>
</div>

<br />

This can lead to a waste of resources, since the idle cores will need to wait for a long time for the busy core(s) to finish its/their huge sequence partitioning (linear time). Therefore, a solution is to not only quick-split, but also "more"-split: in the initial passes, we can partition the input into multiple smaller segments so that even if the splitting itself is more expensive by some constant amount (e.g. due to finer granularity of pivoting), more workers can be started in the subsequent passes. 

<div class="imgcap">
<img src="/QoSort/assets/images/split-in-3.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 5: Quicksort v4.0: more workers in the bootstrap phase helps better the utilization of computing power.</div>
</div>

<br />

We thus introduce a split-by-3 way (**figure 5**), so that we find 2 pivots to make up 3 partitions. At each iteration, the worker will push the two remaining tasks to the task queue and keep working on the largest portion of the 3 partitions. 

This, we found, led to another performance boost, making our quicksort extremely optimized. Sorting with this version of quicksort, the program on 100M elements (each 128 bits) can finish in about 3.9~4.05 seconds, which is quite optimal (and beats the best quicksort and generic-type parallel sorting method that we have found in Go).

##### Miscellaneous

We do want to note that achieving these results was very challenging, since there could be many ways to go. We read some publications on parallel sorting, but most of them were done in GPU, which has limited help for CPU-based parallel sorting. Moreover, we have successfully maintained the in-place sorting invariant so that the memory footprint is kept at a minimal level. This definitely helps with our result!


### QoSort.SampleSort

However, during our analysis, we concluded that quicksort's performance is still bounded by the mechanism through which it works: the \\( n \\)-ary partitioning inevitably leads to relatively compromised locality (because of rescheduling, for example), even though we have made every effort to optimize it within quicksort's bound. Sample sort, on the other hand, is an idea that provides a better solution. 

Our sample sort algorithm is developed base on the algorithm proposed by Professor Guy Blelloch [2], with many strategies adapted for the constraints of GoLang. We decided to proceed with Sample sort because it addresses all flaws with our parallel quick sort algorithm. First, the sampling step ensures input array are partitioned into more balanced sub-arrays, thus work is more balanced. Also upon analysis, we noticed a potential part of improvement for parallel quicksort was at the beginning of the algorithm. At the beginning of quicksort, there aren't many partitions for the worker threads to work on, so we are not able to utilize all the available cores at the beginning of quicksort (this is optimized through the split-by-3, but is still quite suboptimal). 

In sample sort, much of the work is parallelizable, even at the very first iteration of the algorithm. Also sample sort allows us to partition the array into many sub-arrays in the sampling step, which further improves core utilization.

<div class="imgcap">
<img src="/QoSort/assets/images/sample.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 6: Sampling and choosing pivots</div>
</div>

<br />

##### Sampling and Choosing Pivots
In this step, we select a sample set of elements from the input array. The size of the sample set is determined by the number of buckets we plan to use multiplied by a pre-defined over-sampling factor. The elements to be included in the sample set are chosen randomly. After the elements are chosen, we sort the sample set, and choose pivots from the sorted set. The pivots are chosen evenly distributed along the sorted samples, so the pivots have the best chance of dividing the input array into equally sized buckets.


##### Partition Blocks Using Pivots
With appropriate pivots chosen, it remains for us to partition the array into buckets to enable parallel sorting. To do this, we first partitioned the array into equal-sized blocks solely base on indices, then sorted elements within each block using our parallel quick sort algorithm. Since there's no dependency across blocks, Block sorting was carried out in parallel to achieve even better hardware utilization. With each block sorted, we can partition each block into small buckets, using the pivots chosen in the previous step.

<div class="imgcap">
<img src="/QoSort/assets/images/block.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 7: Transposing blocks into buckets</div>
</div>

<br />

##### Transposing Blocks into Buckets

Now we have each block partitioned into small buckets, the next step is to combine small buckets into large buckets and move them to their appropriate location on the array. This step turned out to be very demanding (and interesting). 

In our first attempt, we followed the transpose step implemented in Prof. Blelloch's C++ sample sort, which uses parallel prefix sum (scan). However, when we implemented our version of parallel scan in GoLang, we learned that the overhead for parallel scan is way too much for our input size. In comparison to C++, goroutine for prefix sum is still too heavyweight for frequent spawning and killing. Therefore, after testing, we decided to keep a serial scan. Then, after we obtained prefix sum results, combining buckets and moving data is perfectly suited for parallelism, since there is not interference between different sections of the array.

One optimization that we make to this step is the block transpose. Because we want to transpose `num_blocks x num_buckets` matrix to a `num_buckets x num_blocks` matrix, the eventual transfer of data happen in smaller block matrix. We observed some minor improvements obtained from cache locality by making this optimization. 

##### Sort Each Buckets Using Optimized QuickSort
After the transpose step, we have the array partitioned into similar sized buckets, with all elements in a bucket larger than elements in the previous buckets. All that left to do is to sort the elements within each bucket and we are done. Again, we use our parallel quicksort on each bucket, and the operation across different buckets are also carried out in parallel using a simple go routine and lambda function.


##### Improvement Over QuickSort

Following the implementation specifications above, we were able to address many of the issues we identified with our parallel quicksort algorithm. One problem we had with quicksort was with pivot selection and load balancing. Although we were able to optimize this aspect using split-by-3 technique, there are still space for for improvement. In samplesort, we are able to sample on a large sample set and make more intelligent selection of pivots. In this way, we can partition the array into more smaller chunks of similar size, which is essential in load balancing. Moreover, sample sort is even better suited for adapting parallel computation by design. Most of the operations, especially the block-bucket transpose step, consists of many independent small tasks, which is perfect for utilizing our high core-count CPU architecture. This design also allow us to utilize more cores earlier in the process, as a problem we had with previous iterations is that core utilization is low in early stages of quicksort. 

## IV. Performance and Results

#### Test Setting

Our testing and experiments were carried out on the AFS UNIX machines (Intel Xeon E5-2680), each supplying 20 cores and is 2-way hyperthreaded. The UNIX machines run Golang with version 1.4.2, which provides multi-threading goroutine support.

In order to test the elasticity of our program to different levels of data size, we did the experiments both for 1 million float64-key-value pairs (so 128M bits) as well as 100 million float64-key-value pairs (so 12.8G bits). 

We believe these numbers are reasonable to be counted as "large-scale", since most of the publications on parallel sorting use arrays of sizes well below this line--- even if they use GPU for sorting (e.g. [3]).

#### 1M float64-pairs Sorting

We cross-compare the results of our QoSort library with: (1) Go's built-in sort; and (2) the most popular generic type parallel sort in Go that we found online. Note that this open-source external code for sorting is generic in data type but not data structure: **it assumes inputs of type `[]elem`, which is weaker than our requirement of `sort.Interface`**. 

Here are the results, with performance measured in the time of completing the sort (we take the best out of 8 times):

<div class="imgcap">
<img src="/QoSort/assets/images/1M.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 8: Comparison of sorting performances on 1M float64-pair data.</div>
</div>

<br />

In general, we observed great speedup from existing Go's sorting library (leftmost blue). Moreover, even though the most popular open-source Go sort implementation we found relaxed the generic requirement (see bold face in the paragraph above), it is substantially slower than the optimized quicksort and samplesort. 

#### 100M float64-pairs Sorting

We conduct a similar experiment but with 100M float64-pairs. This is an even more important experiment than the previous one because 100M is truly pushing the boundary of large-scale sorting. The results are shown below:

<div class="imgcap">
<img src="/QoSort/assets/images/100M.png" style="width: 60%">
<div class="thecap" style="color: white; font-size: 14pt">Figure 9: Comparison of sorting performances on 100M float64-pair data.</div>
</div>

<br />

With large-scale data, the improvements from our optimizations are much more obvious. For example, parallel sorts clearly triumphs Go's builtin sorting methods. Moreover, as we optimize the quicksort methods (4 versions as highlighted in orange), eventually we are able to achieve a performance where the program can complete within 3.9 seconds (for split-by-3, and about 4.5 seconds for split-by-2). Sample sort algorithm, moreover, achieved even better result (rightmost bar), only ~2.0 - 2.2 seconds. This is expected, since as discussed above, sample sort offers better parallelization of the tasks, with better resource management (through blocks and buckets) as well as spatial locality (in the transpose phase).

Both our optimized quicksort and samplesort defeated Go's existing most popular sorting extension, which needed around 7.2 seconds for the sorting. Moreover, it is important to note that this open-source implementation relaxed the input type requirement such that it is now an array of generic data type instead of `sort.Interface`.  **Out of curiosity, we tried to relax our requirement to the same degree, and got another 1.5x speedup.**

Based on this result, we consider our implementation relatively successful!

#### Discussion

In terms of large-scale sorting, the speedup is eventually bounded by computation (\\( O(n \log n) \\)). However, given the amount of data, how we deal with locality and scheduling can be really important as well (see the improvement from naive quicksort to optimized quicksort). 

Because of the usage of `sort.Interface` , we cannot exactly measure the relative proportions of time spent in computation and memory access. However, through rough profiling of different regions of our code, the performance is more computation-bounded. To improve on this, a new sorting algorithm with potentially a constant factor fewer comparisons will have to be used. Nevertheless, most of the fast sorting algorithms rely on additional space for purposes like data copying (e.g. mergesort), which leads to higher overhead in data transfer. 

Finally, we would like to re-iterate how important it is to support generic data inputs. Some of the great parallel sorting methods, such as radix sort, rely **heavily** on the underlying element type or data structure (for which array fetch by index is possible). However, we believe it is important to keep the interface contract set by Go's sorting library so that not only ints and floats, but also strings and classes, and be sorted via the same level of abstraction as well.


## V. Acknowledgement 

This has been a rewarding journey where we had the chance to apply the lots of concepts and techniques  learned in 15-418 to real-world usage--- namely, building this library :-) We would like to thank Professor Kayvon for his teaching, and the dedicated TAs for their helps.



## Reference

[1] Median of three partitioning: http://algs4.cs.princeton.edu/23quicksort/

[2] G.E. Blelloch, C.E. Leiserson, B.M. Maggs, C.G. Plaxton, S.J. Smith and M. Zagha. [An Experimental Analysis of Parallel Sorting Algorithms]https://www.cs.cmu.edu/afs/cs.cmu.edu/project/phrensy/pub/papers/BlellochLMPSZ94.pdf] (1998).

[3] D. Božidar and T. Dobravec. [Comparison of parallel sorting algorithms](https://arxiv.org/ftp/arxiv/papers/1511/1511.03404.pdf) (2015).

[4] V. Kale and E. Solomonik. [Parallel Sorting Pattern](http://parlab.eecs.berkeley.edu/wiki/_media/patterns/sortingpattern.pdf) (2010).