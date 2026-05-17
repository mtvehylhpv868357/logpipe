// Package linecount provides a writer wrapper that counts lines written.
package linecount

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counter wraps a writer and tracks how many lines have been written through it.
type Counter struct {
	w     io.WriteCloser
	count int64
}

// New returns a Counter wrapping the given WriteCloser.
func New(w io.WriteCloser) *Counter {
	return &Counter{w: w}
}

// Write passes the line to the underlying writer and increments the counter.
func (c *Counter) Write(line string) error {
	if err := c.w.Write(line); err != nil {
		return err
	}
	atomic.AddInt64(&c.count, 1)
	return nil
}

// Close closes the underlying writer.
func (c *Counter) Close() error {
	return c.w.Close()
}

// Count returns the total number of lines successfully written.
func (c *Counter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Reset resets the line counter to zero.
func (c *Counter) Reset() {
	atomic.StoreInt64(&c.count, 0)
}

// String returns a human-readable summary of lines written.
func (c *Counter) String() string {
	return fmt.Sprintf("lines_written=%d", c.Count())
}
