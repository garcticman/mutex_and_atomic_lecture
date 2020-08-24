package concLesson

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func (f *forecasterShared) isStarted() bool {
	return f.isState(StartState)
}

func (f *forecasterShared) isStopped() bool {
	return f.isState(StopState)
}

func (f *forecasterShared) isState(state uint32) bool {
	return atomic.LoadUint32(f.state) == state
}

func (f *forecasterShared) setStartState() {
	f.setState(StartState)
}

func (f *forecasterShared) setStopState() {
	f.setState(StopState)
}

func (f *forecasterShared) setState(state uint32) {
	atomic.StoreUint32(f.state, state)
}

type forecasterShared struct {
	city   string
	cityMu sync.RWMutex

	state  *uint32
	closed chan struct{}
}

func NewForecasterShared(city string) *forecasterShared {
	return &forecasterShared{
		city:   city,
		state:  new(uint32),
		closed: make(chan struct{}),
	}
}

func (f *forecasterShared) getCity() string {
	f.cityMu.RLock()
	fmt.Printf("RLock %p\n", &f.cityMu)
	defer f.cityMu.RUnlock()

	return f.city
}

func (f *forecasterShared) setCity(city string) {
	f.cityMu.Lock()
	fmt.Printf("Lock %p\n", &f.cityMu)
	defer f.cityMu.Unlock()

	f.city = city
}

// start without waiting
func (f *forecasterShared) start() {
	if !atomic.CompareAndSwapUint32(f.state, StopState, StartState) {
		return
	}

	res := GetWeather(f.city)

	go func() {
		for {
			select {
			case temp := <-res:
				fmt.Printf("Temperature in %s is %d C at %v\n", f.getCity(), temp, time.Now())
			case <-f.closed:
				fmt.Println("Exit")
				return
			}
		}
	}()

	fmt.Println("Started")
}

// stop without waiting
func (f *forecasterShared) stop() {
	if !atomic.CompareAndSwapUint32(f.state, StartState, StopState) {
		return
	}

	close(f.closed)

	fmt.Println("Stopped")
}

func TestCopyMutex(_ *testing.T) {
	city := "London"
	cast := NewForecasterShared(city)

	cast.start()

	go func() {
		timer := time.NewTicker(50 * time.Millisecond)
		defer timer.Stop()
		for range timer.C {
			cast.setCity("Atlanta")
		}
	}()

	time.Sleep(time.Second)

	cast.stop()

	fmt.Println("DONE")
}
