---
layout: post
title: Project Proposal
description: A brief intro to Pastila
---

### Title

Pastila: A Parallel Model of HMM Algorithms

### Authors

  - Shaojie Bai (shaojieb@andrew.cmu.edu)
  - Yutong Chen (yutongc@andrew.cmu.edu)

### Tech

Current plan is develop on C++, enabling potential usage of ISPC and pthreads (or OpenMP or CUDA).

### Summary

This project will focus on parallelizing the **evaluation** and **decoding** problem of Hidden Markov Models (HMM) in large and small state space. While HMM's algorithms (forward, backward & Viterbi) are traditionally iterative, we intend to use a matrix formulation of the problem so as to leverage parallel computing.

### Background

#### What is HMM?

Hidden Markov Model (HMM) is a modeling of a probability distribution. It attempts to describe the domain as a finite set of states \\( S_1, S_2, \dots, S_n \\) as well as a series of functions describing the state transfer. More formally, HMM is a t-tuple \\( (S, \pi, A, B) \\):
  * \\( S \\) is a finite set of states: \\( S_1, \dots, S_n \\). At any time \\( t \\), the state \\(q_t\\) must be one of the states in \\(S\\). 
  * \\( \pi \\) is the prior state probabilities; i.e. "what's the likelihood that initially we are at state \\(S_i\\)?"
  * \\( A \\) is a matrix such that \\( A_{i,j} \\) represents the transition probability from \\( S_i \\) to \\( S_j \\). We call it the **transition matrix**.
  * \\( B \\) is a matrix such that \\( B_{i,j} \\) represents the emission probability of an alphabet \\( V_j \\) at state \\( S_i \\). We call it the **emission matrix**.

<div class="imgcap">
<img src="/images/fg_hmm.png" height="350">
<div class="thecap">Mock transformations of different shapes that (ideally) should drive the distribution to two peaks on two sides of the origin</div>
</div>

Moreover, most of the canonical examples of HMM use symbol \\(O_t\\) to represent the observation at time \\(t\\), and \\(q_t\\) the state at time \\(t\\).
  
In general, the HMM problems are in three areas:
  - Evaluation problem: Compute Probability of observation sequence given a model.
  - Decoding problem: Find state sequence which maximizes probability of observation sequence.  
  - Training problem: Adjust model parameters to maximize probability of observed sequences.

In this project, we plan to focus on the first two problems. With time permitted, we may attempt to tackle parallel training problem as well.
  
#### Forward, Backward and Viterbi Algorithms

There are 3 popular HMM algorithms for evaluation and decoding problems. We represent HMM parameters as \\( \lambda = (\pi, A, B) \\).

##### Forward Algorithm

The forward algorithm tries to determine 
$$
\mathbb{P}[O|\lambda]
$$
(in other words, the probability of observations, given model parameters). To determine this, it essentially uses a DP-style method to compute
$$
\alpha_t(i) = \mathbb{P}[O_1O_2 \dots O_t, q_t=S_i | \lambda]
$$
which is essentially the probability that, given parameters \\( \lambda \\), we observed \\(O_1, \dots, O_t\\) and stay at state \\(S_i\\) at time \\(t\\). 

Since transition is modeled using matrix \\(A\\), we can use this to compute \\(\alpha_{t+1}\\) and so on. Eventually, we should have 
$$
\mathbb{P}[O|\lambda] = \sum_{i} \alpha_T(i)
$$
where \\(T\\) is the total time (i.e. final time).

 
##### Backward Algorithm
 
Similar idea as the forward algorithm, except that we proceed backwards (from \\(T\\) to \\(0\\)). 

##### Viterbi Algorithm

At time \\(t\\), the Viterbi Algorithm is essentially trying to determine
$$
\delta_t(i) = \max_{q_1, \dots, q_{t-1}} \mathbb{P}[q_1, \dots, q_{t-1}, q_t=S_i, O_1, O_2, \dots, O_t | \lambda]
$$

#### The Challenge

There are three aspects in this problem that we found interesting:
  - **Inherent sequential nature of the algorithm**: These algorithms are traditionally implemented in an iterative way. This is not surprising since the algorithm itself relies on dynamic-programming, which goes from time \\(t\\) (e.g. \\( \alpha_t \\)) to time \\(t+1\\).
  - **Large size of data**: Moreover, when the # of states get larger, memory may come into picture, and it has been found that the evaluation and decoding problem can take time measured in days and weeks. Some data sets, such as genome sequencing and speech recognition, have large observations, which makes efficient analysis of the HMM even harder.
  - **Insufficient parallel works**: We didn't find any useful open-source programs available online that supports parallel HMM algorithms we try to implement here.

#### Resources
  
  - **Hardware**: We need multicore machines that support SIMD vector programming. If time permitted, we may also try CUDA version. **It'd be great if the course staff can help us with this!!**
  - **References**: 
    * L. Rabiner, "A tutorial on hidden Markov models and selected applications in speech recognition," pp.267-296, 1990
    * J. Nielsen & A. Sand, "Algorithms for a parallel implementation of Hidden Markov Models with a small state space," IEEE International Parallel & Distributed Processing Symposium, 2011
    * A. Sand, "Engineering of Algorithms for Hidden Markov Models and Tree Distances", PhD Dissertation, 2014
    
We may append to this list as the project goes.   
   

### Goals

#### What do we want to achieve?

##### Speedup
Considering the iterative nature of the algorithm, we expect to achieve a speedup of at least 2x compared to the sequential algorithm. **Note that we are not employing the trivially parallel computing here; we don't simply send in multiple streams of input to evaluate**. Instead, we shall take a method that focus more on the problem formulation (i.e. the theoretical) itself.

If things went well, we certainly will try to hit a better mark. Sand reports to be able to achieve a 4x speedup under some special circumstances (small states, small data set), but we certainly will try to generalize this, if we are beyond the schedule!

##### Demo Plan (Deliverables)

We expect to present the speedup graphs compared to the most popular sequential model. Moreover, we hope to be able to be able to demo the evaluation and decoding on some real dataset, which we will look for in kaggle (after we finish the implementation). 

The C++ code shall also be available online. 
 
##### Baseline
 
We will implement a sequential model as a baseline to compare to first.
 

### Platform choice

We will develop mainly in C++. Its support of pthreads, vector programs (SIMD), CUDA and OpenMP offers us maximum flexibility in terms of researching.

### Tentative schedule

  - By April 15: Finish the baseline implementation as well as necessary paper reading.
  - By April 20: Develop necessary structures and codes needed for later customizations.
  - By April 25: Developed some ideas for parallelizing the tasks, and start implementing them.
  - By May 1: Finish at least one parallel modeling of forward and Viterbi algorithms.
  - By May 8: Explore potential improvements based on the problems that we run into, and further optimize the model(s). Run on real dataset if time permitted. 
  - By May 10: Have all necessary data collected, and be prepared for the presentation.
   
  