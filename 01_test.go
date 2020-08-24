package concLesson

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// race example
func TestWithoutMutex(t *testing.T) {
	done := make(chan bool, 3)

	counter := 0
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			for counter = 0; counter < 5; counter++ {
				fmt.Printf("goroutine: %d, counter = %d \n", i, counter)
			}
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestSyncMutex(t *testing.T) {
	mutex := sync.Mutex{}
	done := make(chan bool, 3)

	counter := 0
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			mutex.Lock()
			for counter = 0; counter < 5; counter++ {
				fmt.Printf("goroutine: %d, counter = %d \n", i, counter)
			}
			mutex.Unlock()
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestAtomic(t *testing.T) {
	done := make(chan bool, 3)

	counter := uint64(0)
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			for atomic.StoreUint64(&counter, 0); atomic.LoadUint64(&counter) < 5; atomic.AddUint64(&counter, 1) {
				fmt.Printf("goroutine: %d, counter = %d \n", i, atomic.LoadUint64(&counter))
			}
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}
