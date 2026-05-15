// Package dedupe provides line deduplication for log streams.
// It suppresses repeated identical lines within a configurable window.
package dedupe

import (
	"sync"
	"time"
)

// Filter suppresses duplicate lines seen within a sliding time window.
type Filter struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	clock   func() time.Time
}

// New creates a Filter that deduplicates lines seen within window duration.
// A zero or negative window disables deduplication (all lines pass through).
func New(window time.Duration) *Filter {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Filter {
	return &Filter{
		seen:   make(map[string]time.Time),
		window: window,
		clock:  clock,
	}
}

// Allow returns true if the line should be forwarded.
// A line is suppressed if it was seen within the configured window.
func (f *Filter) Allow(line string) bool {
	if f.window <= 0 {
		return true
	}

	now := f.clock()

	f.mu.Lock()
	defer f.mu.Unlock()

	f.evict(now)

	if _, exists := f.seen[line]; exists {
		return false
	}

	f.seen[line] = now
	return true
}

// evict removes entries that have expired. Must be called with mu held.
func (f *Filter) evict(now time.Time) {
	for line, ts := range f.seen {
		if now.Sub(ts) >= f.window {
			delete(f.seen, line)
		}
	}
}

// Reset clears all tracked lines immediately.
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]time.Time)
}
