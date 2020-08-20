package main

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestAtomicUint(t *testing.T) {
	var n uint64

	atomic.AddUint64(&n, 100)
	fmt.Println(n)
	atomic.AddUint64(&n, ^uint64(10 - 1))
	fmt.Println(n)
	old := atomic.SwapUint64(&n, 100)
	fmt.Println(n, old)

	swapped := atomic.CompareAndSwapUint64(&n, 100, 90)
	fmt.Println(swapped, n)
	swapped = atomic.CompareAndSwapUint64(&n, 100, 100)
	fmt.Println(swapped, n)
}
