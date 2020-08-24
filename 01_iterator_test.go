package concLesson

import (
	"fmt"
	"testing"
)

func TestIterator(t *testing.T) {
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			fmt.Printf("goroutine: %d\n", i)
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestIteratorFixed(t *testing.T) {
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		i := i
		go func() {
			fmt.Printf("goroutine: %d\n", i)
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}
