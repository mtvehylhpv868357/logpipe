package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counters tracks line processing statistics across the pipeline.
type Counters struct {
	LinesRead    atomic.Int64
	LinesMatched atomic.Int64
	LinesDropped atomic.Int64
	WriteErrors  atomic.Int64
}

// IncRead increments the lines-read counter.
func (c *Counters) IncRead() {
	c.LinesRead.Add(1)
}

// IncMatched increments the lines-matched counter.
func (c *Counters) IncMatched() {
	c.LinesMatched.Add(1)
}

// IncDropped increments the lines-dropped counter.
func (c *Counters) IncDropped() {
	c.LinesDropped.Add(1)
}

// IncWriteError increments the write-error counter.
func (c *Counters) IncWriteError() {
	c.WriteErrors.Add(1)
}

// Snapshot returns a point-in-time copy of the current counter values.
type Snapshot struct {
	LinesRead    int64
	LinesMatched int64
	LinesDropped int64
	WriteErrors  int64
}

// Snapshot captures the current counter values atomically.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:    c.LinesRead.Load(),
		LinesMatched: c.LinesMatched.Load(),
		LinesDropped: c.LinesDropped.Load(),
		WriteErrors:  c.WriteErrors.Load(),
	}
}

// Print writes a human-readable summary of the snapshot to w.
func (s Snapshot) Print(w io.Writer) {
	fmt.Fprintf(w, "lines read:    %d\n", s.LinesRead)
	fmt.Fprintf(w, "lines matched: %d\n", s.LinesMatched)
	fmt.Fprintf(w, "lines dropped: %d\n", s.LinesDropped)
	fmt.Fprintf(w, "write errors:  %d\n", s.WriteErrors)
}
