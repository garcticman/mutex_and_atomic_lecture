package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDeferUnlock(t *testing.T) {
	mutex := sync.Mutex{}
	value := 0

	for i := 0; i < 100; i++ {
		i := i
		go func() {
			mutex.Lock()
			defer mutex.Unlock()

			value += i

			returnFunc := func() {
				fmt.Println(value)
				return
			}

			if value > 4000 {
				returnFunc()
			}
			if i > 50 && i < 60 {
				returnFunc()
			}
			if value == i {
				returnFunc()
			}
			if value > 10 && value < 20 {
				value += i
				return
			}
		}()
	}

	time.Sleep(time.Second * 2)
}
