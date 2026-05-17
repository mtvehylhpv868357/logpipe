// Package tail provides a writer wrapper that keeps only the last N lines
// written to it, discarding earlier lines once the limit is reached.
package tail

import (
	"fmt"
	"io"
)

// Limiter is a writer that forwards lines to an underlying writer but only
// after it has buffered enough lines to satisfy the tail window. When the
// window is full the oldest line is evicted before the new one is accepted.
// A zero or negative max disables the limiter and every line is forwarded
// immediately.
type Limiter struct {
	w      io.WriteCloser
	max    int
	buf    []string
	head   int
	count  int
}

// New returns a Limiter that keeps the last max lines. When the Limiter is
// closed all buffered lines are flushed to w in order.
func New(w io.WriteCloser, max int) *Limiter {
	if max <= 0 {
		return &Limiter{w: w}
	}
	return &Limiter{
		w:   w,
		max: max,
		buf: make([]string, max),
	}
}

// Write accepts a single log line. When max is zero the line is forwarded
// immediately. Otherwise it is stored in the ring buffer; if the buffer was
// already full the oldest entry is silently dropped.
func (l *Limiter) Write(line string) error {
	if l.max <= 0 {
		return l.w.Write(line)
	}
	l.buf[l.head] = line
	l.head = (l.head + 1) % l.max
	if l.count < l.max {
		l.count++
	}
	return nil
}

// Close flushes all buffered lines in chronological order to the underlying
// writer and then closes it.
func (l *Limiter) Close() error {
	if l.max <= 0 {
		return l.w.Close()
	}
	start := 0
	if l.count == l.max {
		start = l.head
	}
	for i := 0; i < l.count; i++ {
		idx := (start + i) % l.max
		if err := l.w.Write(l.buf[idx]); err != nil {
			return fmt.Errorf("tail: flush: %w", err)
		}
	}
	return l.w.Close()
}
