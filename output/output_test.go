package output

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// bufWriter is a test helper that captures written lines.
type bufWriter struct {
	buf    bytes.Buffer
	closed bool
}

func (b *bufWriter) Write(line string) error {
	fmt.Fprintln(&b.buf, line)
	return nil
}
func (b *bufWriter) Close() error { b.closed = true; return nil }

func TestStdoutWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	w := &StdoutWriter{out: &buf}
	if err := w.Write("hello stdout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello stdout") {
		t.Errorf("expected output to contain 'hello stdout', got %q", buf.String())
	}
}

func TestFileWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	fw, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("NewFileWriter: %v", err)
	}
	if err := fw.Write("file line 1"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := fw.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "file line 1") {
		t.Errorf("expected file to contain 'file line 1', got %q", string(data))
	}
}

func TestFileWriter_InvalidPath(t *testing.T) {
	_, err := NewFileWriter("/nonexistent/dir/file.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestMulti_Write(t *testing.T) {
	a, b := &bufWriter{}, &bufWriter{}
	m := NewMulti(a, b)

	if err := m.Write("broadcast"); err != nil {
		t.Fatalf("Multi.Write: %v", err)
	}
	for _, w := range []*bufWriter{a, b} {
		if !strings.Contains(w.buf.String(), "broadcast") {
			t.Errorf("expected writer to contain 'broadcast', got %q", w.buf.String())
		}
	}
}

func TestMulti_Close(t *testing.T) {
	a, b := &bufWriter{}, &bufWriter{}
	m := NewMulti(a, b)
	if err := m.Close(); err != nil {
		t.Fatalf("Multi.Close: %v", err)
	}
	if !a.closed || !b.closed {
		t.Error("expected all writers to be closed")
	}
}
