package qosortv2

import (
    "math/rand"
    "sync"
	"math"
)


var QUICKSORT_THRESHOLD = 10000
var OVER_SAMPLE = 8

// set to 0 will spawn num_cores worker threads in each call to Qsort_par
var QSORT_POOL_SIZE = 0

func SampleSort(A []qselem) {
    sample_sort(A, 0, len(A))
}

func sample_sort(A []qselem, i, j int) {
    n := j - i
    if n < QUICKSORT_THRESHOLD {
        Qsort_parallel(A, i, j, 0)
        return
    }

    // TODO: NEED TO FIND BEST BLOCK/BUCKET NUMBERS
    num_blocks := int(math.Ceil(math.Sqrt(float64(n)) / 4))
    block_size := ((n-1)/num_blocks) + 1
    num_buckets := num_blocks
    sample_set_size := num_buckets * OVER_SAMPLE
    m := num_blocks * num_buckets
    
    sample_set := make([]qselem, sample_set_size)

    // Randomly sample from input
    // parallel?
    for i := 0; i < sample_set_size; i++ {
        s := rand.Int() % n
        sample_set[i] = A[s]
    }

    // TODO: Test whether we should use qsort_serial
    Qsort_serial(sample_set, 0, sample_set_size)

    // evenly select pivots from sorted sample
    pivots := make([]qselem, num_buckets-1)
    for k := 0; k < num_buckets-1; k++ {
        pivots[k] = sample_set[OVER_SAMPLE * k + OVER_SAMPLE/2]
    }

    sketch := make([]qselem, n)
    counts := make([]int, m)
    copy(sketch, A)

    // sort within each block and count size of each bucket
    wg := new(sync.WaitGroup)
    for b := 0; b < num_blocks; b++ {
		b_copy := b
        wg.Add(1)
        go func() {
            offset := b_copy * block_size
            size := block_size
            if b_copy == num_blocks - 1 { size = n - offset }      // The last block will take whatever's left
            Qsort_parallel(sketch, offset, offset + size, QSORT_POOL_SIZE)
            merge_seq(sketch, offset, size, pivots, num_buckets-1, counts, b_copy*num_buckets)
            wg.Done()
        }()
    }
    wg.Wait()

	bucket_offsets := transpose_buckets(sketch, A, counts, n, block_size, num_blocks, num_buckets)

    wg2 := new(sync.WaitGroup)
    for b := 0; b < num_buckets; b++ {
		b_copy := b
        wg2.Add(1)
        go func() {
            istart := bucket_offsets[b_copy]
            iend := bucket_offsets[b_copy+1]
            Qsort_parallel(A, istart, iend, QSORT_POOL_SIZE)
            wg2.Done()
        }()
    }
    wg2.Wait()
}

/*
 * Parameters:
 *     A: Copy of the original array (sketch).
 *     A_offset: The start index of the block.
 *     A_size: The size of the data this block is responsible for (starting from A_offset)
 *     pivots: The bucket separator pivots; shared by all blocks.
 *     num_pivots: Total # of pivots.
 *     counts: Array of integer counters for each bucket within each block.
 *     count_offset: The start index in counts where this block will need to access.
 */
func merge_seq(A []qselem, A_offset int, A_size int,
               pivots []qselem, num_pivots int,
               counts []int, count_offset int) {
    ia := A_offset
    ib := 0
    ic := count_offset
    if A_size == 0 || num_pivots == 0 {return}   // Unlikely to happen
    for i := 0; i <= num_pivots; i++ {counts[count_offset + i] = 0}
    for {
        for A[ia].Less(pivots[ib]) {
            counts[ic]++
            ia++
            if ia == A_offset + A_size {return}
        }
		// Invariant: A[ia] >= pivots[ib]
        ib++
        ic++
        if ib == num_pivots {break}

        // If pivots[ib-1] == pivots[ib], do duplicate check.
        if !pivots[ib-1].Less(pivots[ib]) {
			// While A[ia] == pivots[ib] == pivots[ib-1], put it in the duplicate range.
            for !pivots[ib].Less(A[ia]) {
                counts[ic]++
                ia++
                if ia == A_offset + A_size {return}
            }
            ib++
            ic++
            if ib == num_pivots {break}
        }
    }
    counts[ic] = A_offset + A_size - ia

}
