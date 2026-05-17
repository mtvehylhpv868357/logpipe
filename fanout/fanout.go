// Package fanout provides a writer that duplicates each line to multiple writers.
package fanout

import "io"

// Writer writes each line to all underlying writers.
// If a writer fails, the error is recorded but writing continues to remaining writers.
type Writer struct {
	writers []io.WriteCloser
}

// New returns a Writer that fans out to all provided writers.
func New(writers ...io.WriteCloser) *Writer {
	return &Writer{writers: writers}
}

// Write writes p to every underlying writer.
// Returns the first error encountered, if any, but always attempts all writers.
func (f *Writer) Write(p []byte) (int, error) {
	var firstErr error
	for _, w := range f.writers {
		if _, err := w.Write(p); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return len(p), firstErr
}

// WriteString writes s to every underlying writer.
func (f *Writer) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

// Close closes all underlying writers.
// Returns the first error encountered, but always attempts to close all.
func (f *Writer) Close() error {
	var firstErr error
	for _, w := range f.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Len returns the number of underlying writers.
func (f *Writer) Len() int {
	return len(f.writers)
}
