package linecount_test

import (
	"errors"
	"sync"
	"testing"

	"logpipe/linecount"
)

// mockWriter is a simple WriteCloser used in tests.
type mockWriter struct {
	lines []string
	failWrite bool
}

func (m *mockWriter) Write(line string) error {
	if m.failWrite {
		return errors.New("write error")
	}
	m.lines = append(m.lines, line)
	return nil
}

func (m *mockWriter) Close() error { return nil }

func TestCounter_InitialCountIsZero(t *testing.T) {
	c := linecount.New(&mockWriter{})
	if c.Count() != 0 {
		t.Fatalf("expected 0, got %d", c.Count())
	}
}

func TestCounter_IncrementsOnWrite(t *testing.T) {
	c := linecount.New(&mockWriter{})
	_ = c.Write("hello")
	_ = c.Write("world")
	if c.Count() != 2 {
		t.Fatalf("expected 2, got %d", c.Count())
	}
}

func TestCounter_DoesNotIncrementOnError(t *testing.T) {
	c := linecount.New(&mockWriter{failWrite: true})
	err := c.Write("fail")
	if err == nil {
		t.Fatal("expected error")
	}
	if c.Count() != 0 {
		t.Fatalf("expected 0, got %d", c.Count())
	}
}

func TestCounter_Reset(t *testing.T) {
	c := linecount.New(&mockWriter{})
	_ = c.Write("a")
	_ = c.Write("b")
	c.Reset()
	if c.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", c.Count())
	}
}

func TestCounter_String(t *testing.T) {
	c := linecount.New(&mockWriter{})
	_ = c.Write("x")
	s := c.String()
	if s != "lines_written=1" {
		t.Fatalf("unexpected string: %s", s)
	}
}

func TestCounter_ConcurrentWrites(t *testing.T) {
	c := linecount.New(&mockWriter{})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.Write("line")
		}()
	}
	wg.Wait()
	if c.Count() != 100 {
		t.Fatalf("expected 100, got %d", c.Count())
	}
}
