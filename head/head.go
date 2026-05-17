// Package head implements a line limiter that passes through only the first N lines.
package head

import "sync/atomic"

// Limiter passes through at most Max lines, then drops all subsequent lines.
type Limiter struct {
	max   int64
	count atomic.Int64
}

// New returns a Limiter that allows at most max lines through.
// If max is zero or negative, the limiter is disabled and all lines pass through.
func New(max int64) *Limiter {
	return &Limiter{max: max}
}

// Allow returns true if the line should be forwarded.
func (l *Limiter) Allow() bool {
	if l.max <= 0 {
		return true
	}
	n := l.count.Add(1)
	return n <= l.max
}

// Reset resets the internal counter so lines are allowed again from the start.
func (l *Limiter) Reset() {
	l.count.Store(0)
}

// Remaining returns how many more lines will be allowed through.
// Returns -1 if the limiter is disabled.
func (l *Limiter) Remaining() int64 {
	if l.max <= 0 {
		return -1
	}
	r := l.max - l.count.Load()
	if r < 0 {
		return 0
	}
	return r
}
