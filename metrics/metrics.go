package metrics

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// Counters holds atomic counters for pipeline telemetry.
type Counters struct {
	received uint64
	routed   uint64
	dropped  uint64
	errors   uint64
}

// NewCounters returns a zeroed Counters instance.
func NewCounters() *Counters {
	return &Counters{}
}

func (c *Counters) IncReceived() { atomic.AddUint64(&c.received, 1) }
func (c *Counters) IncRouted()   { atomic.AddUint64(&c.routed, 1) }
func (c *Counters) IncDropped()  { atomic.AddUint64(&c.dropped, 1) }
func (c *Counters) IncErrors()   { atomic.AddUint64(&c.errors, 1) }

// Snapshot is an immutable copy of counter values at a point in time.
type Snapshot struct {
	Received uint64
	Routed   uint64
	Dropped  uint64
	Errors   uint64
}

// Snapshot returns a consistent read of all counters.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		Received: atomic.LoadUint64(&c.received),
		Routed:   atomic.LoadUint64(&c.routed),
		Dropped:  atomic.LoadUint64(&c.dropped),
		Errors:   atomic.LoadUint64(&c.errors),
	}
}

// Print writes a human-readable summary to w (defaults to os.Stdout).
func (s Snapshot) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "received=%d routed=%d dropped=%d errors=%d\n",
		s.Received, s.Routed, s.Dropped, s.Errors)
}

// DropRate returns the fraction of received messages that were dropped,
// in the range [0.0, 1.0]. Returns 0 if no messages have been received.
func (s Snapshot) DropRate() float64 {
	if s.Received == 0 {
		return 0
	}
	return float64(s.Dropped) / float64(s.Received)
}
