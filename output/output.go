package output

import (
	"fmt"
	"io"
	"os"
)

// Writer is the interface that all output destinations must implement.
type Writer interface {
	Write(line string) error
	Close() error
}

// StdoutWriter writes log lines to stdout.
type StdoutWriter struct {
	out io.Writer
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{out: os.Stdout}
}

func (w *StdoutWriter) Write(line string) error {
	_, err := fmt.Fprintln(w.out, line)
	return err
}

func (w *StdoutWriter) Close() error { return nil }

// FileWriter appends log lines to a file.
type FileWriter struct {
	path string
	file *os.File
}

func NewFileWriter(path string) (*FileWriter, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("output: open file %q: %w", path, err)
	}
	return &FileWriter{path: path, file: f}, nil
}

func (w *FileWriter) Write(line string) error {
	_, err := fmt.Fprintln(w.file, line)
	return err
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}

// Multi fans a log line out to multiple Writers.
type Multi struct {
	writers []Writer
}

func NewMulti(writers ...Writer) *Multi {
	return &Multi{writers: writers}
}

func (m *Multi) Write(line string) error {
	var firstErr error
	for _, w := range m.writers {
		if err := w.Write(line); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (m *Multi) Close() error {
	var firstErr error
	for _, w := range m.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
