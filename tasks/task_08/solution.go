package main

import (
	"math"
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	var now time.Time
	if clock != nil {
		now = clock.Now()
	} else {
		now = time.Now() // just for initialization to pass
	}

	limiter := Limiter{
		mu:     sync.Mutex{},
		clock:  clock,
		rate:   math.Max(0, ratePerSec),
		burst:  burst,
		tokens: float64(burst),
		last:   now}
	return &limiter
}

func (l *Limiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.burst <= 0 || l.clock == nil || n <= 0 || n > l.burst {
		return false
	}

	now := l.clock.Now()
	since := now.Sub(l.last)
	if since.Seconds() > 0 {
		l.last = now
	}

	// add tokens since last call
	// shouldn't exceed l.burst and be less than zero
	l.tokens = math.Min(float64(l.burst), math.Max(0, l.tokens+since.Seconds()*l.rate))

	if nf := float64(n); l.tokens >= nf {
		l.tokens -= nf
		return true
	} else {
		return false
	}
}
