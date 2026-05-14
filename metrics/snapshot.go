package metrics

import (
	"fmt"
	"io"
	"time"
)

// Snapshot holds a point-in-time copy of counter values.
type Snapshot struct {
	Timestamp time.Time
	Received  int64
	Routed    int64
	Dropped   int64
	Errors    int64
}

// TakeSnapshot returns a Snapshot from the current Counters state.
func (c *Counters) TakeSnapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		Timestamp: time.Now(),
		Received:  c.received,
		Routed:    c.routed,
		Dropped:   c.dropped,
		Errors:    c.errors,
	}
}

// Print writes a human-readable summary of the snapshot to w.
func (s Snapshot) Print(w io.Writer) {
	fmt.Fprintf(w, "[%s] received=%d routed=%d dropped=%d errors=%d\n",
		s.Timestamp.Format(time.RFC3339),
		s.Received,
		s.Routed,
		s.Dropped,
		s.Errors,
	)
}
