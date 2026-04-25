package breaker

import (
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
