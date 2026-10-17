package contingency

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
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	lastFailTime time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     StateClosed,
		threshold: threshold,
		cooldown:  cooldown,
	}
}

func (cb *CircuitBreaker) ShouldUseContingency() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.cooldown {
			cb.state = StateHalfOpen
			return false
		}
		return true
	}
	return false
}

func (cb *CircuitBreaker) ReportResult(isInfraFailure bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if !isInfraFailure {
		cb.failures = 0
		cb.state = StateClosed
		return
	}

	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.threshold || cb.state == StateHalfOpen {
		cb.state = StateOpen
	}
}
