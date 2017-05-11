---
layout: post
title: Checkpoint Report
description: Achievements so far...
image: assets/images/brain.jpg
---

## Schedule Ahead

  - By the end of April, we shall have a high-performance sorting method implemented in Go with generic type support.
  - By May 4, we shall have at least 2 (or 3, if possible) parallel sorting methods supported in Go so that we are ready for the competition. We will focus on quicksort and sample sort for the competition.
  - By May 6, we will try to further optimize the Go implementations.
  - By May 10, we should be ready for the presentation.
 

## Completed Work

We started with rewriting the quicksort algorithm written in C++ provided by Prof. Guy Blelloch into GoLang. Our initial version was designed to only work on ``float64``, for which we were able to achieve comparable performance to GoLang's native generic sort. However, when we tried to adapt the float sort to a generic-typed sort (using ``sort.Interface``), we experienced significant performance bump because of the use of ``Interface`` feature. Upon some research and testing, it seems there's no good solution to this problem currently, so we had to improve the performance in algorithm tweaking and parallelism scheduling.

We also analyzed and optimized the parallelism so that the method better work in GoLang. In particular, the C++ implementation used a split-by-3 method, which we found not to work very well in Go due to GoLang's requirement of implementing a ``sort.Interface`` for the data structure. Moreover, the fact that Go manages multithreading in a different way makes the underlying parallelization different. Instead, we focused on the more traditional 2-way-split. We sampled the points so that the pivot we pick each time is much closer to the median. Moreover, we apply careful synchronization on the goroutines (i.e. threads in Go) so that in-place sorting can work perfectly--- this helps avoid the need to additional memory.

## Goals and Deliverables
__For 210-sorting competition (due May 4)__: Since the deadline for 210-sorting competition much sooner than the 418 project deadline, the plan for the following week is to focus on developing a fast parallel sort in GoLang, and tuning to the given test machine and test dataset, so we can have a competitive algorithm for the competition.

__For 418 project__: The focus of 418 project submission would be a detailed analysis on performance of parallel sorting across various sorting algorithms and tuning. We will adapt from the sorting algorithm we implemented for the 210-sorting competition, and test the performance of algorithm in different settings. Plus, we plan to test performance of common parallel sorting algorithms including quicksort, samplesort, and bitonic sort, implemented in GoLang. As a stretch goal, we may perform some testing in Java and compare the result to GoLang.


## What to Show at the Parallelism Competition
For the parallelism competition, we plan to present graphs and anlysis comparing different parallel sorting techniques, and perhaps across different programing languages.


## Preliminary Result
We tested the different sorting methods we implemented in GoLang so far on Unix5 machine, which has 2 10-core Intel Xeon E5-2680 v2 processors, resulting a total of 40 logical cores. The results are presented below.

#### Serial Quicksort

We managed to implement a quicksort (with carefully picked pivot) that supports generic type and is **faster** than Go's builtin sorting method ``sort.Sort``. Testing on a random float64-pair array of size 10M, `qosort.qsort_serial` is able to complete the task within 4.6 seconds (on average), whereas the builtin sort needs about 5.6 seconds. This is a 1.2x speedup. 

#### Parallel Quicksort

Of the various parallel quicksort implementations we offered in our library (3-way split, 2-way split, etc.), our code is currently able to sort 10M double-pair (128 bits per element) in 0.675 second **on a machine with 20 cores**. The competition will be on a 72-core machine, so we believe this is a surprisingly good performance, especially given that we can further optimize it. Note that for ordinary quicksort with perfectly chosen pivor, the parallel wall-clock time complexity is \\( O(\log^2 n) \\). 

## Issues/Concerns/Unknowns

Naturally, we have encountered quite a few challenges/obstacles so far.

  - **Generic Type**: As discussed above, our intention to support generic sorting in Go has created some performance issue. For example, we found that for the identically implemented insertion sort (one supports generic typing while the other one only `float64`), the former one takes longer time in general. Therefore, while we are eager to have better performance for this competition, we also ought to consider the usability of this library in general.
  - **Bad parallel framework**: Unlike C++, GoLang does not have a mature parallel framework or library that supports multi-core execution. The only place where this can be specified is Go's `runtime` package. While there are open-source repos that provide SIMD support for Go, the core logic was still written in the Assembly and such packages are forbidden in the competition in general.
  - **Time**: We shifted to this project from a previous one during the Spring break. Therefore, the time for us to catch up was very limited.
  - **Hardware**: Currently, we are testing on the Unix5 machine of the AFS. However, the competition offers the opportunity of execution on a 72-core machine. Therefore, we truly hope that we can get access to a server of similar property (e.g. Intel machine, etc.).


