package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of lines per time window.
type Limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	count    int
	windowAt time.Time
	now      func() time.Time
}

// New creates a Limiter that allows at most maxLines lines per window duration.
// A maxLines value <= 0 disables rate limiting (all lines pass).
func New(maxLines int, window time.Duration) *Limiter {
	return &Limiter{
		max:      maxLines,
		window:   window,
		windowAt: time.Now(),
		now:      time.Now,
	}
}

// newWithClock is used in tests to inject a custom clock.
func newWithClock(maxLines int, window time.Duration, clock func() time.Time) *Limiter {
	l := New(maxLines, window)
	l.now = clock
	l.windowAt = clock()
	return l
}

// Allow returns true if the line should be forwarded under the current rate limit.
func (l *Limiter) Allow() bool {
	if l.max <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if now.Sub(l.windowAt) >= l.window {
		l.count = 0
		l.windowAt = now
	}

	if l.count < l.max {
		l.count++
		return true
	}
	return false
}

// Reset resets the counter and starts a new window immediately.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.count = 0
	l.windowAt = l.now()
}
