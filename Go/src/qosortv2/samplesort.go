package qosortv2

import (
    "math/rand"
    "sync"
)


QUICKSORT_THRESHOLD := 20000
OVER_SAMPLE := 8

func SampleSort(A []qselem) {
    sample_sort(A, 0, len(A))
}

func sample_sort(A []qselem, i, j int) {
    n := j - i
    if (n < QUICKSORT_THRESHOLD) {
        Qsort_parallel(A, i, j)
        return
    }

    // NEED TO FIND BEST BLOCK/BUCKET NUMBERS
    num_blocks := 6
    block_size := ((n-1)/num_blocks) + 1;
    num_buckets := 6
    sample_set_size := num_buckets * OVER_SAMPLE;
    m := num_blocks * num_buckets
    
    sample_set := new([]qselem, sample_set_size)

    // Randomly sample from input
    // parallel?
    for i := 0; i < sample_set_size; i++ {
        s := rand.Int() % n
        sample_set[i] = A[s]
    }

    // should we use qsort_serial ???
    Qsort_serial(sample_set, 0, sample_set_size)

    // evenly select pivots from sorted sample
    pivots = new([]qselem, num_buckets-1)
    for k := 0; k < num_buckets-1; k++ {
        pivots[k] = sample_set[OVER_SAMPLE * k + OVER_SAMPLE/2]
    }

    sketch := new([]qselem, n)
    counts := new([]int, m)
    copy(sketch, A)

    // sort within each block and count size of each bucket
    wg := new(sync.WaitGroup)
    for b := 0; b < num_blocks; b++ {
        wg.Add(1)
        go func() {
            offset := b * block_size;
            size = (i < num_blocks - 1) ? block_size : n-offset;
            Qsort_parallel(sketch, offset, size)
            merge_seq(sketch, offset, size, pivots, num_buckets-1, counts, b*num_buckets)
            wg.Done()
        }()
    }
    wg.Wait()

    // **************************
    // add transpose buckets here
    // **************************

    wg2 := new(sync.WaitGroup)
    for b := 0; b < num_buckets; b++ {
        wg2.Add(1)
        go func() {
            istart := bucket_offsets[b]
            iend := bucket_offsets[b+1]
            Qsort_parallel(A, istart, iend)
            wg.Done()
        }()
    }
    wg2.Wait()
}

func merge_seq(A []qselem, A_offset int, A_size int,
               pivots []int, num_pivots int,
               counts []int, count_offset int) {
    ia = A_offset
    ib = 0
    ic = count_offset
    if A_size == 0 || num_pivots == 0 {return}
    for i := 0; i <= num_pivots; i++ {counts[count_offset + i] = 0}
    for {
        for A[ia].Less(pivots[ib]) {
            counts[ic]++
            ia++
            if ia == A_offset + A_size {return}
        }
        ib++
        ic++
        if ib == num_pivots {break}

        if !pivots[ib-1].Less(pivots[ib]) {
            for !pivots[ib].Less(A[ia]) {
                counts[ic]++
                ia++
                if ia == A_offset + A_size {return}
            }
            ib++
            ic++
            if (ib == num_pivots) {break}
        }
    }
    counts[ic] = A_offset + A_size - ia

}
