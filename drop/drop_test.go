package drop

import (
	"errors"
	"testing"
)

// mockWriter records written lines and optionally returns an error.
type mockWriter struct {
	lines []string
	fail  bool
}

func (m *mockWriter) Write(line string) error {
	if m.fail {
		return errors.New("write error")
	}
	m.lines = append(m.lines, line)
	return nil
}

func (m *mockWriter) Close() error { return nil }

func TestDropper_AllowsNonMatchingLine(t *testing.T) {
	next := &mockWriter{}
	d, err := New(`ERROR`, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Write("INFO hello world"); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if len(next.lines) != 1 || next.lines[0] != "INFO hello world" {
		t.Errorf("expected line to be forwarded, got %v", next.lines)
	}
}

func TestDropper_DropsMatchingLine(t *testing.T) {
	next := &mockWriter{}
	d, err := New(`ERROR`, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Write("ERROR something went wrong"); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if len(next.lines) != 0 {
		t.Errorf("expected line to be dropped, got %v", next.lines)
	}
}

func TestDropper_EmptyPatternForwardsAll(t *testing.T) {
	next := &mockWriter{}
	d, err := New("", next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = d.Write("anything")
	if len(next.lines) != 1 {
		t.Errorf("expected line forwarded, got %v", next.lines)
	}
}

func TestDropper_InvalidPatternReturnsError(t *testing.T) {
	next := &mockWriter{}
	_, err := New(`[invalid`, next)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestDropper_ForwardsWriteError(t *testing.T) {
	next := &mockWriter{fail: true}
	d, err := New(`SKIP`, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Write("INFO line"); err == nil {
		t.Fatal("expected write error to be propagated")
	}
}

func TestDropper_Close(t *testing.T) {
	next := &mockWriter{}
	d, _ := New(``, next)
	if err := d.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}
