// Package sampling provides log line sampling strategies for logpipe.
// It allows reducing output volume by forwarding only a fraction of matching lines.
package sampling

import (
	"math/rand"
	"sync/atomic"
)

// Sampler decides whether a given log line should be forwarded.
type Sampler interface {
	Sample(line string) bool
}

// RateSampler forwards lines with a given probability in [0.0, 1.0].
type RateSampler struct {
	rate float64
	rng  *rand.Rand
}

// NewRateSampler returns a RateSampler that forwards lines with the given rate.
// rate must be between 0.0 and 1.0; values outside this range are clamped.
func NewRateSampler(rate float64, src rand.Source) *RateSampler {
	if rate < 0.0 {
		rate = 0.0
	}
	if rate > 1.0 {
		rate = 1.0
	}
	return &RateSampler{rate: rate, rng: rand.New(src)}
}

// Sample returns true with probability equal to the configured rate.
func (s *RateSampler) Sample(_ string) bool {
	return s.rng.Float64() < s.rate
}

// NthSampler forwards every Nth line and drops the rest.
type NthSampler struct {
	n       uint64
	counter atomic.Uint64
}

// NewNthSampler returns a NthSampler that forwards every nth line.
// n must be >= 1; a value of 1 forwards all lines.
func NewNthSampler(n uint64) *NthSampler {
	if n == 0 {
		n = 1
	}
	return &NthSampler{n: n}
}

// Sample returns true for every nth call.
func (s *NthSampler) Sample(_ string) bool {
	v := s.counter.Add(1)
	return v%s.n == 0
}

// PassthroughSampler forwards all lines (rate = 1.0).
type PassthroughSampler struct{}

// Sample always returns true.
func (PassthroughSampler) Sample(_ string) bool { return true }
