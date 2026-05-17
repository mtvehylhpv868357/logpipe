// Package timestamp provides a writer wrapper that prepends a timestamp to each line.
package timestamp

import (
	"fmt"
	"io"
	"time"
)

const defaultFormat = time.RFC3339

// clockFn allows injection of a custom clock for testing.
type clockFn func() time.Time

// Stamper wraps an io.WriteCloser and prepends a timestamp to each written line.
type Stamper struct {
	out    io.WriteCloser
	format string
	clock  clockFn
}

// New returns a Stamper using RFC3339 timestamps.
func New(out io.WriteCloser) *Stamper {
	return NewWithFormat(out, defaultFormat)
}

// NewWithFormat returns a Stamper using the given time format string.
func NewWithFormat(out io.WriteCloser, format string) *Stamper {
	if format == "" {
		format = defaultFormat
	}
	return &Stamper{
		out:    out,
		format: format,
		clock:  time.Now,
	}
}

func newWithClock(out io.WriteCloser, format string, fn clockFn) *Stamper {
	s := NewWithFormat(out, format)
	s.clock = fn
	return s
}

// Write prepends the current timestamp to line and writes to the underlying writer.
func (s *Stamper) Write(line string) error {
	stamped := fmt.Sprintf("%s %s", s.clock().Format(s.format), line)
	return s.out.Write(stamped)
}

// Close closes the underlying writer.
func (s *Stamper) Close() error {
	return s.out.Close()
}
