package concLesson

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
			defer func() {
				if v := recover(); v != nil {
					fmt.Println(v)
				}
				mutex.Unlock()
			}()

			value += i

			if value > 4000 {
				fmt.Println(value)
				return
			}
			if i > 50 && i < 60 {
				fmt.Println(value)
				return
			}
			if value == i {
				fmt.Println(value)
				return
			}
			if value > 10 && value < 20 {
				value += i
				return
			}
			if i == 10 {
				panic("Hello")
			}
		}()
	}

	time.Sleep(time.Second * 2)
}

func TestPanicDeadlock(t *testing.T) {
	mutex := sync.Mutex{}
	var value uint64

	inc := func() {
		mutex.Lock()
		value++
		mutex.Unlock()
	}
	incWithPanic := func() {
		mutex.Lock()
		value++
		panic("Hello")
		mutex.Unlock()
	}

	func() {
		defer func() {
			v := recover()
			fmt.Println("v =", v)
		}()
		incWithPanic()
	}()
	fmt.Println("before inc deadlock")
	inc()
}
