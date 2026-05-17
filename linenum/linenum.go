// Package linenum provides a transformer that prepends a line number to each line.
package linenum

import (
	"fmt"
	"sync/atomic"
)

// Numberer prepends an incrementing line number to each line.
type Numberer struct {
	counter int64
	format  string
}

// New returns a Numberer using the default format "[%d] ".
func New() *Numberer {
	return NewWithFormat("[%d] ")
}

// NewWithFormat returns a Numberer that uses the given fmt format string.
// The format must contain exactly one integer verb (e.g. "%d" or "%05d").
func NewWithFormat(format string) *Numberer {
	if format == "" {
		format = "[%d] "
	}
	return &Numberer{format: format}
}

// Transform prepends the next line number to line and returns the result.
func (n *Numberer) Transform(line string) string {
	num := atomic.AddInt64(&n.counter, 1)
	return fmt.Sprintf(n.format, num) + line
}

// Reset sets the internal counter back to zero.
func (n *Numberer) Reset() {
	atomic.StoreInt64(&n.counter, 0)
}

// Count returns the number of lines transformed so far.
func (n *Numberer) Count() int64 {
	return atomic.LoadInt64(&n.counter)
}
