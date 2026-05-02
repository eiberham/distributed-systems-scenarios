package breaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	sync.Mutex
	name             string
	failureThreshold int
	state            State
	failureCount     int
	lastFailureTime  time.Time
	resetTimeout     time.Duration
}

func New(name string, failureThreshold int, state State, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:             name,
		failureThreshold: failureThreshold,
		state:            state,
		resetTimeout:     resetTimeout,
	}
}

func (cb *CircuitBreaker) Run(call func() (interface{}, error)) (interface{}, error) {
	cb.Mutex.Lock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
		} else {
			cb.Mutex.Unlock()
			return nil, errors.New("circuit is open")
		}
	}
	cb.Mutex.Unlock()

	response, err := call()

	cb.Mutex.Lock()
	defer cb.Mutex.Unlock()

	if err != nil {
		cb.recordFailure()
		return nil, err
	}

	cb.recordSuccess()
	return response, nil
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.failureCount = 0
	if cb.state == StateHalfOpen {
		cb.changeState(StateClosed)
	}
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failureCount++

	// If we're in half open and fail, or if we hit the threshold in closed
	if cb.state == StateHalfOpen || cb.failureCount >= cb.failureThreshold {
		cb.changeState(StateOpen)
	}
}

func (cb *CircuitBreaker) changeState(newState State) {
	cb.state = newState
	if newState == StateOpen {
		cb.lastFailureTime = time.Now()
	}
}
