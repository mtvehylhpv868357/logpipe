// Package ceiling provides a writer wrapper that enforces a maximum byte
// size per line, dropping lines that exceed the limit.
package ceiling

import "fmt"

// Writer drops any line whose byte length exceeds MaxBytes.
// When MaxBytes is zero or negative the limiter is disabled and every
// line is forwarded unchanged.
type Writer struct {
	dst      LineWriter
	maxBytes int
}

// LineWriter is the interface satisfied by any downstream writer.
type LineWriter interface {
	Write(line string) error
}

// New returns a Writer that forwards lines to dst only when their byte
// length is within maxBytes. Set maxBytes <= 0 to disable the ceiling.
func New(dst LineWriter, maxBytes int) *Writer {
	return &Writer{dst: dst, maxBytes: maxBytes}
}

// Write forwards line to the downstream writer when the line is within
// the configured byte ceiling. Lines that exceed the ceiling are silently
// dropped and a descriptive error is returned so callers can track drops.
func (w *Writer) Write(line string) error {
	if w.maxBytes <= 0 {
		return w.dst.Write(line)
	}
	if len(line) > w.maxBytes {
		return fmt.Errorf("ceiling: line length %d exceeds max %d, dropped", len(line), w.maxBytes)
	}
	return w.dst.Write(line)
}
