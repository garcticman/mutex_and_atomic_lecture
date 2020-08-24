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
	atomic.AddUint64(&n, ^uint64(10-1))
	fmt.Println(n)
	old := atomic.SwapUint64(&n, 100)
	fmt.Println(n, old)

	swapped := atomic.CompareAndSwapUint64(&n, 100, 90)
	fmt.Println(swapped, n)
	notSwapped := atomic.CompareAndSwapUint64(&n, 100, 100)
	fmt.Println(notSwapped, n)
}

func TestArbitraryType(t *testing.T) {
	type arbitraryType struct {
		a, b, c int
	}
	value := arbitraryType{1, 2, 3}

	var atomicValue atomic.Value
	atomicValue.Store(value)
	loaded := atomicValue.Load().(arbitraryType)
	fmt.Println(loaded)
	fmt.Println(value == loaded)
}
