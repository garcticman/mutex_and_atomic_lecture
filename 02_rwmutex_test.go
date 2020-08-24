package concLesson

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	goroutineCount = 10
	doneLen        = goroutineCount * 2
)

func BenchmarkMutexes(b *testing.B) {
	b.Run("RW", RWMutex)
	b.Run("Simple", SimpleMutex)
}

type someObject struct {
	sync.Mutex
	someState map[string]uint64
}

func newObject() someObject {
	return someObject{
		someState: make(map[string]uint64),
	}
}

func SimpleMutex(b *testing.B) {
	example := newObject()

	done := make(chan bool, doneLen)

	for i := uint64(0); i < goroutineCount; i++ {
		i := i
		go func() {
			example.Lock()
			example.someState["name"+strconv.FormatUint(i, 10)] = i
			example.Unlock()
			done <- true
		}()
	}

	for i := uint64(0); i < goroutineCount; i++ {
		i := i
		go func() {
			example.Lock()

			time.Sleep(time.Millisecond)
			_ = example.someState["name"+strconv.FormatUint(i, 10)]

			example.Unlock()
			done <- true
		}()
	}
	for i := 0; i < doneLen; i++ {
		<-done
	}
}

type RWObject struct {
	sync.RWMutex
	someState map[string]uint64
}

func NewRWObject() RWObject {
	return RWObject{
		someState: make(map[string]uint64),
	}
}

func RWMutex(b *testing.B) {
	example := NewRWObject()

	done := make(chan bool, doneLen)

	for i := uint64(0); i < goroutineCount; i++ {
		i := i
		go func() {

			example.Lock()
			example.someState["name"+strconv.FormatUint(i, 10)] = i
			example.Unlock()
			done <- true
		}()
	}

	for i := uint64(0); i < goroutineCount; i++ {
		i := i
		go func() {
			example.RLock()

			time.Sleep(time.Millisecond)
			_ = example.someState["name"+strconv.FormatUint(i, 10)]

			example.RUnlock()
			done <- true
		}()
	}
	for i := 0; i < doneLen; i++ {
		<-done
	}
}

func (sm *RWObject) SetValue(key string, value uint64) {
	sm.Lock()
	sm.someState[key] = value
	sm.Unlock()
}

func (sm *RWObject) GetValue(key string) uint64 {
	sm.RLock()
	time.Sleep(time.Millisecond)
	value := sm.someState[key]
	sm.RUnlock()

	return value
}

func TestSimpleMutexWithoutDancing(t *testing.T) {
	example := NewRWObject()

	done := make(chan bool, doneLen)

	for i := uint64(0); i < 10; i++ {
		i := i
		go func() {
			example.SetValue("name"+strconv.FormatUint(i, 10), i)

			done <- true
		}()
	}

	for i := uint64(0); i < goroutineCount; i++ {
		i := i
		go func() {
			_ = example.GetValue("name" + strconv.FormatUint(i, 10))

			done <- true
		}()
	}
	for i := 0; i < doneLen; i++ {
		<-done
	}
}
