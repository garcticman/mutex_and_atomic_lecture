package concLesson

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type forecaster struct {
	city   string
	state  *uint32
	closed chan struct{}

	sync.Once

	sync.Mutex
}

const (
	StopState = iota
	StartState
)

func newForecaster(city string) *forecaster {
	return &forecaster{city, new(uint32), make(chan struct{}, 1), sync.Once{}, sync.Mutex{}}
}

func (f *forecaster) isStarted() bool {
	return f.isState(StartState)
}

func (f *forecaster) isStopped() bool {
	return f.isState(StopState)
}

func (f *forecaster) isState(state uint32) bool {
	return atomic.LoadUint32(f.state) == state
}

func (f *forecaster) setStartState() {
	f.setState(StartState)
}

func (f *forecaster) setStopState() {
	f.setState(StopState)
}

func (f *forecaster) setState(state uint32) {
	atomic.StoreUint32(f.state, state)
}

// start without waiting
func (f *forecaster) start() {
	f.Do(func() {
		fmt.Println("Was initiated. Now we can start.")
	})

	/*
		// fixme it's a trap
		if currentState := atomic.LoadUint32(f.state); currentState != StopState {
			return
		} else {
			atomic.StoreUint32(f.state, StartState)
		}
	*/

	/*
		// possible but ...
		f.Mutex.Lock()
		if *f.state == StartState {
			f.Mutex.Unlock()

			return
		}
		f.Mutex.Unlock()
	*/

	if !atomic.CompareAndSwapUint32(f.state, StopState, StartState) {
		return
	}

	res := GetWeather(f.city)

	go func() {
		for {
			select {
			case temp := <-res:
				fmt.Printf("Temperature in %s is %d C at %v\n", f.city, temp, time.Now())
			case <-f.closed:
				fmt.Println("Exit")
				return
			}
		}
	}()

	fmt.Println("Started")
}

// stop without waiting
func (f *forecaster) stop() {
	/*
		f.Mutex.Lock()
		if *f.state == StopState {
			f.Mutex.Unlock()

			return
		}
		f.Mutex.Unlock()
	*/

	if !atomic.CompareAndSwapUint32(f.state, StartState, StopState) {
		return
	}

	close(f.closed)

	fmt.Println("Stopped")
}

func TestStartStopStates(_ *testing.T) {
	city := "London"
	cast := newForecaster(city)

	go cast.start()
	go cast.start()
	go cast.start()
	go cast.start()

	time.Sleep(10 * time.Second)

	go cast.stop()
	go cast.stop()
	go cast.stop()
	go cast.stop()

	fmt.Println("DONE")
}

func GetWeather(_ string) chan int {
	res := make(chan int, 1)

	go func() {
		for {
			res <- GetWeatherByTime("", time.Now())
		}
	}()

	return res
}

// GetWeather gets a temperature forecast from remote API by the town name and time
func GetWeatherByTime(_ string, _ time.Time) int {
	res := int(rand.Int31n(45))
	switch {
	case res%2 == 0:
		time.Sleep(500 * time.Millisecond)
	case res%3 == 0:
		time.Sleep(800 * time.Millisecond)
	case res%5 == 0:
		time.Sleep(3 * time.Second)
	case res%10 == 0:
		time.Sleep(time.Minute)
	default:
		time.Sleep(300 * time.Millisecond)
	}
	return res
}
