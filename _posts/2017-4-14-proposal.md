---
layout: post
title: Project Proposal
description: A brief intro to Pastila
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

This project is inspired by CMU'S [210 Sorting Competition](http://www.cs.cmu.edu/~15210/sort.html), where we are to design and implement sorting algorithms that can work well on large-scale datasets (of size ~100M). We will explore two different programming languages (Java and Go), along with different sorting methods, to build a simple parallel sorting library.

## Background

The two major factors that this project shall focus on are (1) language differences (GC, in particular); and (2) sorting methods (how parallel is it?)

#### 1. Garbage-Collected Languages

TBD (shaojieb)
  
#### 2. Sample Sort

TBD (yutongc)

#### 3. Radix Sort

TBD (shaojieb)


## The Challenge

There are three aspects in this problem that we found interesting:
  - **Large scale**: With the data size on a scale of 100M (i.e. 100,000,000), we expect memory management to be an important issue. How to best exploit the cache, with the garbage collection differences in mind, is the core question.
  - **Task parallelization**: How to parallelize the steps so that we can best utilize the CPUs as well as vector lanes (if we are to use it)? This can depend on the sorting method we are using. Moreover, we need to pay attention to the potential communication cost involved.

## Resources
  
  - **Hardware**: We need multicore machines that support SIMD vector programming. If time permitted, we may also try CUDA version. **It'd be great if the course staff can help us with this!!**
  - **References**: 
    * TBA
    * TBA
    * TBA
    
We may append to this list as the project goes.   
   

## Goals

### What do we want to achieve?

#### 1. Speedup

For the competition, Professor Guy Blleloch provided runnable codes in SML and C++ that used **sample sort** method. In particular, the SML code was able to sort a size-100M array in 0.33 second on a 72-core machine and 1.1 seconds on a 20-core machine. The C++ implementation (with essentially the same idea) could achieve the same goal in 0.183 second on a 72-core machine.

Our goal, of

#### 2. Demo Plan (Deliverables)

We expect to present the speedup graphs compared to the most popular sequential model. Moreover, we hope to be able to be able to demo the evaluation and decoding on some real dataset, which we will look for in kaggle (after we finish the implementation). 

The C++ code shall also be available online. 
 
We will also implement a sequential model as a baseline to compare to first.
 

## Platform choice

We will develop mainly in C++. Its support of pthreads, vector programs (SIMD), CUDA and OpenMP offers us maximum flexibility in terms of researching.

## Tentative schedule

  - By April 15: Finish the baseline implementation as well as necessary paper reading.
  - By April 20: Develop necessary structures and codes needed for later customizations.
  - By April 25: Developed some ideas for parallelizing the tasks, and start implementing them.
  - By May 1: Finish at least one parallel modeling of forward and Viterbi algorithms.
  - By May 8: Explore potential improvements based on the problems that we run into, and further optimize the model(s). Run on real dataset if time permitted. 
  - By May 10: Have all necessary data collected, and be prepared for the presentation.
   
  