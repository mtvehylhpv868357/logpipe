// Package tee provides a writer that copies each line to a secondary writer
// while passing it through to the primary writer unchanged.
package tee

import "io"

// Writer wraps a primary io.WriteCloser and copies every written line to a
// secondary io.Writer (e.g. os.Stderr for debug inspection).
type Writer struct {
	primary   io.WriteCloser
	secondary io.Writer
}

// New returns a Writer that writes each line to primary and also copies it to
// secondary. If secondary is nil the writer behaves as a transparent pass-through.
func New(primary io.WriteCloser, secondary io.Writer) *Writer {
	return &Writer{primary: primary, secondary: secondary}
}

// Write writes p to the primary writer and, if a secondary writer is set,
// copies p to the secondary writer. Errors from the secondary writer are
// silently discarded so that a failing side-channel never interrupts the
// primary pipeline.
func (w *Writer) Write(p []byte) (int, error) {
	if w.secondary != nil {
		_, _ = w.secondary.Write(p)
	}
	return w.primary.Write(p)
}

// Close closes the primary writer. The secondary writer is not closed because
// it may be shared (e.g. os.Stderr).
func (w *Writer) Close() error {
	return w.primary.Close()
}
