package ceiling_test

import (
	"errors"
	"strings"
	"testing"

	"logpipe/ceiling"
)

// mockWriter collects written lines and can optionally return an error.
type mockWriter struct {
	lines []string
	err   error
}

func (m *mockWriter) Write(line string) error {
	if m.err != nil {
		return m.err
	}
	m.lines = append(m.lines, line)
	return nil
}

func TestCeiling_DisabledWhenZero(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, 0)
	long := strings.Repeat("x", 200)
	if err := w.Write(long); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(dst.lines) != 1 {
		t.Fatalf("expected 1 line forwarded, got %d", len(dst.lines))
	}
}

func TestCeiling_DisabledWhenNegative(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, -10)
	if err := w.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst.lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(dst.lines))
	}
}

func TestCeiling_AllowsLineWithinLimit(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, 10)
	if err := w.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst.lines) != 1 || dst.lines[0] != "hello" {
		t.Fatalf("expected line to be forwarded, got %v", dst.lines)
	}
}

func TestCeiling_AllowsLineAtExactLimit(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, 5)
	if err := w.Write("abcde"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst.lines) != 1 {
		t.Fatalf("expected line forwarded at exact limit")
	}
}

func TestCeiling_DropsLineExceedingLimit(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, 5)
	err := w.Write("toolongline")
	if err == nil {
		t.Fatal("expected error for oversized line, got nil")
	}
	if len(dst.lines) != 0 {
		t.Fatalf("expected no lines forwarded, got %d", len(dst.lines))
	}
}

func TestCeiling_ErrorContainsLengths(t *testing.T) {
	dst := &mockWriter{}
	w := ceiling.New(dst, 4)
	err := w.Write("hello world")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "11") || !strings.Contains(err.Error(), "4") {
		t.Errorf("error should mention lengths, got: %v", err)
	}
}

func TestCeiling_PropagatesDstError(t *testing.T) {
	sentinel := errors.New("write failed")
	dst := &mockWriter{err: sentinel}
	w := ceiling.New(dst, 100)
	if err := w.Write("ok"); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
