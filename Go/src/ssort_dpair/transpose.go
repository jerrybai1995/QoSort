package ssort_dpair

import "sync"

var TRANS_THRESHOLD = 64


/************************* Transpose class *************************/
type Transpose struct { A, B []int }

func (T Transpose) transR(rStart int, rCount int, rLength int,
		                  cStart int, cCount int, cLength int) {
	if cCount < TRANS_THRESHOLD && rCount < TRANS_THRESHOLD {
		// TODO: Should have been parallel for
		for i := rStart; i < rStart + rCount; i++ {
			for j := cStart; j < cStart + cCount; j++ {
				T.B[j*cLength + i] = T.A[i*rLength + j]
			}
		}
	} else if cCount > rCount {
		l1 := cCount / 2
		l2 := cCount - l1
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			T.transR(rStart, rCount, rLength, cStart, l1, cLength)
			wg.Done()
		}()
		T.transR(rStart, rCount, rLength, cStart + l1, l2, cLength)
		wg.Wait()
	} else {
		l1 := rCount/2
		l2 := rCount - l1
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			T.transR(rStart, l1, rLength, cStart, cCount, cLength)
			wg.Done()
		}()
		T.transR(rStart + l1, l2, rLength, cStart, cCount, cLength)
		wg.Wait()
	}
}

func (T Transpose) Trans(rCount int, cCount int) {
	T.transR(0, rCount, cCount, 0, cCount, rCount)
}


/************************* BlockTrans class *************************/
type BlockTrans struct {
	A, B []doublepair
	OA, OB, L []int
}

func (T BlockTrans) transR(rStart int, rCount int, rLength int,
	                       cStart int, cCount int, cLength int) {
	if cCount < TRANS_THRESHOLD && rCount < TRANS_THRESHOLD {
		// TODO: Should have been parallel for
		for i := rStart; i < rStart + rCount; i++ {
			for j := cStart; j < cStart + cCount; j++ {
				pa := T.OA[i*rLength + j]
				pb := T.OB[j*cLength + i]
				l := T.L[i*rLength + j]
				for k := 0; k < l; k++ {
					T.B[pb + k] = T.A[pa + k]
				}
			}
		}
	} else if cCount > rCount {
		l1 := cCount / 2
		l2 := cCount - l1
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			T.transR(rStart, rCount, rLength, cStart, l1, cLength)
			wg.Done()
		}()
		T.transR(rStart, rCount, rLength, cStart + l1, l2, cLength)
		wg.Wait()
	} else {
		l1 := rCount / 2
		l2 := rCount - l1
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			T.transR(rStart, l1, rLength, cStart, cCount, cLength)
			wg.Done()
		}()
		T.transR(rStart + l1, l2, rLength, cStart, cCount, cLength)
		wg.Wait()
	}
}

func (T BlockTrans) Trans(rCount int, cCount int) {
	T.transR(0, rCount, cCount, 0, cCount, rCount)
}


/************************* transpose_buckets impl*************************/

func scan_add_seq(A []int, output []int) {
	n := len(A)
	cumsum := 0
	for i := 0; i < n; i++ {
		temp := A[i]
		output[i] = cumsum
		cumsum += temp
	}
}

func scan_add_par(A []int, n int, output []int) {
	N := n
	copy(output, A)
	wg := new(sync.WaitGroup)
	for twod := 1; twod < N; twod *= 2 {
		twod1 := twod*2
		for i := 0; i < N; i += twod1 {
			j := i
			wg.Add(1)
			go func() {
				output[j+twod1-1] += output[j+twod-1]
				wg.Done()
			}()
		}
		wg.Wait()
	}

	output[N-1] = 0;

	for twod := N/2; twod >= 1; twod /= 2 {
		twod1 := twod*2
		for i := 0; i < N; i += twod1 {
			j := i
			wg.Add(1)
			go func() {
				t := output[j+twod-1]
				output[j+twod-1] = output[j+twod1-1]
				output[j+twod1-1] += t
				wg.Done()
			}()
		}
		wg.Wait()
	}


}

func transpose_buckets(from []doublepair, to []doublepair, counts []int,
                       n int, block_size int, num_blocks int,
                       num_buckets int) []int {
	m := num_buckets * num_blocks
	dest_offsets := make([]int, m)

	source_offsets := make([]int, m)
	seq_counts := make([]int, m)
	copy(seq_counts, counts)

	scan_add_seq(seq_counts, source_offsets)

	Transpose{counts, dest_offsets}.Trans(num_blocks, num_buckets)

	scan_add_seq(dest_offsets, dest_offsets)

	BlockTrans{from, to, source_offsets, dest_offsets, counts}.Trans(num_blocks, num_buckets)

	bucket_offsets := make([]int, num_buckets+1)

	for i := 0; i < num_buckets; i++ { bucket_offsets[i] = dest_offsets[i * num_blocks] }
	bucket_offsets[num_buckets] = n

	return bucket_offsets
}