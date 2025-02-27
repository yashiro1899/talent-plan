package main

import (
	"runtime"
	"sync"
)

var tokens chan struct{}

func init() {
	tokens = make(chan struct{}, runtime.NumCPU()*2)
}

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	n := len(src)
	if n < 2 {
		return
	}

	sort64(src)

	// 开辟一个与原来数组一样大小的空间用来存储用
	temp := make([]int64, n)
	for i := 64; i < n; i *= 2 {
		var wg sync.WaitGroup
		for left, right := 0, 0; left < n-i; left = right {
			mid := left + i
			right = mid + i
			if right > n {
				right = n
			}

			tokens <- struct{}{} // acquire a token
			wg.Add(1)
			go merge(src, temp[left:right], left, mid, mid, right, &wg)
		}
		wg.Wait()
	}
}

func sort64(src []int64) {
	n := len(src)
	for start := 0; start < n; start += 64 {
		end := start + 64
		if end > n {
			end = n
		}

		// expect := src[start:end]
		// sort.Slice(expect, func(i, j int) bool { return expect[i] < expect[j] })
		for i := start + 1; i < end; i++ {
			for j := i; j > start && src[j] < src[j-1]; j-- {
				src[j], src[j-1] = src[j-1], src[j]
			}
		}
	}
}

func merge(src, section []int64, a, b, c, d int, wg *sync.WaitGroup) {
	next := 0
	for a < b && c < d {
		if src[a] < src[c] {
			section[next] = src[a]
			a++
		} else {
			section[next] = src[c]
			c++
		}
		next++
	}

	// 上面循环结束的条件有两个，如果是左边的游标尚未到达，那么需要把
	// 数组接回去，可能会有疑问，那如果右边的没到达呢，其实模拟一下就可以
	// 知道，如果右边没到达，那么说明右边的数据比较大，这时也就不用移动位置了
	if b-a > 0 {
		e := c - (b - a)
		copy(src[e:c], src[a:b])
		c = e
	}
	copy(src[c-next:c], section[:next])

	<-tokens // release the token
	wg.Done()
}
