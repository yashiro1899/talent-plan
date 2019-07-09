package main

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	n := len(src)
	if n < 2 {
		return
	}

	temp := make([]int64, n)

	for i := 1; i < n; i *= 2 {
		var a, b, c, d int

		for a = 0; a < n-i; a = d {
			b = a + i
			c = b
			d = b + i
			if d > n {
				d = n
			}

			next := 0
			for a < b && c < d {
				if src[a] < src[c] {
					temp[next] = src[a]
					a++
				} else {
					temp[next] = src[c]
					c++
				}
				next++
			}

			for a < b {
				c--
				b--
				src[c] = src[b]
			}

			for next > 0 {
				c--
				next--
				src[c] = temp[next]
			}
		}
	}
}
