// Package throttle provides a line-based throughput limiter that delays
// writes to enforce a maximum lines-per-second rate across an output.
package throttle

import (
	"time"
)

// Throttler enforces a maximum number of lines per second by sleeping
// between writes when the rate would be exceeded.
type Throttler struct {
	maxPerSec int
	interval  time.Duration
	sleep     func(time.Duration)

	windowStart time.Time
	count       int
}

// New creates a Throttler that allows at most maxPerSec lines per second.
// If maxPerSec is zero or negative, the Throttler is disabled (pass-through).
func New(maxPerSec int) *Throttler {
	return newWithSleep(maxPerSec, time.Sleep)
}

func newWithSleep(maxPerSec int, sleep func(time.Duration)) *Throttler {
	var interval time.Duration
	if maxPerSec > 0 {
		interval = time.Second / time.Duration(maxPerSec)
	}
	return &Throttler{
		maxPerSec:   maxPerSec,
		interval:    interval,
		sleep:       sleep,
		windowStart: time.Now(),
	}
}

// Allow blocks until the current line may be forwarded according to the
// configured rate. It returns immediately when the Throttler is disabled.
func (t *Throttler) Allow() {
	if t.maxPerSec <= 0 {
		return
	}

	now := time.Now()

	// Reset window every second.
	if now.Sub(t.windowStart) >= time.Second {
		t.windowStart = now
		t.count = 0
	}

	t.count++

	if t.count > t.maxPerSec {
		// Overshoot: sleep for one interval to smooth the rate.
		t.sleep(t.interval)
		return
	}

	// Pace writes evenly within the window.
	expected := t.windowStart.Add(time.Duration(t.count-1) * t.interval)
	if delay := expected.Sub(now); delay > 0 {
		t.sleep(delay)
	}
}

// Disabled reports whether the Throttler is a no-op.
func (t *Throttler) Disabled() bool {
	return t.maxPerSec <= 0
}
