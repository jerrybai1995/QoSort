---
layout: post
title: Project Proposal
description: A brief intro to QoSort
image: assets/images/ai.jpg
---

## Title

QuSort: Optimized Large-Scale Parallel Sorting on CPUs

## Authors

  - Shaojie Bai (shaojieb@andrew.cmu.edu)
  - Yutong Chen (yutongc@andrew.cmu.edu)

## Tech

  - **Programming Languages**: Java, Go
  - **Hardware Tech**: CPU (w/ OpenMP), potentially SIMD

## Summary

This project is inspired by CMU'S [210 Sorting Competition](http://www.cs.cmu.edu/~15210/sort.html), where we are to design and implement sorting algorithms that can work well on large-scale datasets (of size at magnitude ~100M). However, our goal is more than just **a fastest algorithm**, but also a study of parallelizing sorting tasks optimally, in order to build a simple parallel sorting library.

## Background

The two major factors that this project shall focus on are:
  * Comparing and analyzing different ways to parallelize a selected set of sorting methods
  * How to support generic sorting interfaces in a garbage collected language (in this case, GoLang). 

#### 1. Garbage-Collected Languages

Unlike languages such as C or C++, many languages support garbage collection through mechanisms such as reference counting. Such automatic mechanism signals a tradeoff of development convenience and execution performance--- especially with respect to memory management and computation resources. In this project, we focus on the sorting library in Go, which is a garbage collected language.
  
#### 2. Quick Sort

Quicksort is a classical sorting method. However, while easy to do this iteratively (and in an embarrassingly parallel way), it is significantly harder to optimize the parallelism within it. Given a huge array of double values, load balance, cache management and sorting method itself will have enormous effect on the performance. 

#### 3. Sample Sort

One idea that we have in mind is the parallel sample sort. This is similar to quicksort, but the idea is now to randomly pick a few bucket separator. Then, parallel sorting can take place within each bucket. 


## The Challenge

There are three aspects in this problem that we found interesting:
  - **No existing great parallel sorting library in the language**: The existing sort in Go, for example, is strictly sequential (albeit fast). Therefore, we expect to spend lots of time spotting the parallelism within the algorithm, and meanwhile create a benchmark to compare to.
  - **Most work done on GPU**: We found that many of the peer-evaluated publications have been focusing on GPU-based large-scale sorting. For those that address CPU-based parallel sorting, the authors typically use C or C++--- so they are not very helpful in our case.
  - **Generic interface support**: GoLang's builtin sort supports not only array, but also generic sorting interface `sort.Interface`, where only three methods are implemented: `Len()`, `Swap(i, j)` and `Less(i, j)`. This means no indexing (i.e. getting and setting), which restricts the potential sorting methods we can leverage.
  - **High dependency on sorting method**: Different sorting methods create lots of differences. Some methods may not be sequential-friendly, but can be great to use in parallel environment (e.g. bitonic sort). Therefore, it is our job to explore different methods and analyze their pros and cons.

## Resources
  
  - **Hardware**: We need multicore (at least 24) machines that (potentially) support SIMD vector programming. **It'd be great if the course staff can help us with this!!**
  - **References**: 
    * D. Božidar and T. Dobravec. [Comparison of parallel sorting algorithms](https://arxiv.org/ftp/arxiv/papers/1511/1511.03404.pdf) (2015).
    * P. Tsigas and Y. Zhang. [A Simple, Fast Parallel Implementation of Quicksort and its Performance Evaluation on SUN Enterprise 10000](http://www.cse.chalmers.se/~tsigas/papers/Pquick.pdf) (2003).
    * TBA
    
We may append to this list as the project goes.   
   

## Goals

### What do we want to achieve?

#### 1. Speedup

For the competition, Professor Guy Blleloch provided runnable codes in SML and C++ that used **sample sort** method. In particular, the SML code was able to sort a size-100M array in 0.33 second on a 72-core machine and 1.1 seconds on a 20-core machine. The C++ implementation (with essentially the same idea) could achieve the same goal in 0.183 second on a 72-core machine.

Our goal, of course, is to implement a parallel sort that is reasonably fast. Because we are primarily using a higher-level language, we expect the performance to be not as good as the program in C++. But in general, we hope to get within 3 or 4 seconds for sorting on array of size 100M.

#### 2. Demo Plan (Deliverables)

We expect to demo the comparison graphs, where we shall compare the results across different sorting methods. We will also discuss in detail what ideas we've tried to optimize the tasks.
 

## Platform choice

We will develop mainly in Go. While it may be worthwhile to try other languages as well, focusing on a good well-supported Go library is equally (if not more) important! 

## Tentative schedule

  - By April 18: Finish the baseline implementation as well as necessary paper reading. Identify benchmark. 
  - By April 20: Develop the first parallel version of quicksort and identify areas of improvement.
  - By April 25: Develop some ideas for further optimizing parallelism on the tasks; start on sample sort. 
  - By April 27: At least finish with one sorting model. Ought to make the sorting interface as generic as possible (e.g. not only limiter to []float64 in Go, etc.).
  - By May 3: Finish at least two sorting models in Go. Be prepared for the competition.
  - By May 10: Have all necessary data collected, and be prepared for the presentation.
   
  