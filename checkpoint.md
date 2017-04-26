## Schedule


## Summarize completed work

(1-2 paragraphs)

We started with rewriting the quicksort algorithm written in c++ provided by Prof. Guy Blelloch into GoLang. Our initial version was designed to only work on ``float64``, for which we were able to achieve comparable performance to GoLang's native generic sort. However, when we tried to adapt the float sort to a generic-typed sort (using ``sort.Interface``), we experienced significant performance bump because of the use of ``Interface`` feature. Upon some research and testing, it seems there's no good solution to this problem currently, so we had to improve the performance in algorithm tweaking and parallelism scheduling.



## Goals and deliverables
__For 210-sorting competition (due May 4)__: Since the deadline for 210-sorting competition much sooner than the 418 project deadline, the plan for the following week is to focus on developing a fast parallel sort in GoLang, and tuning to the given test machine and test dataset, so we can have a competitive algorithm for the competition.

__For 418 project__: The focus of 418 project submission would be a detailed analysis on performance of parallel sorting across various sorting algorithms and tuning. We will adapt from the sorting algorithm we implemented for the 210-sorting competition, and test the performance of algorithm in different settings. Plus, we plan to test performance of common parallel sorting algorithms including quicksort, samplesort, and bitonic sort, implemented in GoLang. As a stretch goal, we plan to perform the same testing in Java and compare the result to GoLang.


## What to show at the parallelism competition
For the parallelism competition, we plan to present graphs and anlysis comparing different parallel sorting techniques, and perhaps across different programing languages.


## preliminary result
We tested the different sorting methods we implemented in GoLang so far on Unix5 machine, which has 2 10-core Intel Xeon E5-2680 v2 processors, resulting a total of 40 logical cores. The results are presented below.

## Issues/Concerns/Unknowns


