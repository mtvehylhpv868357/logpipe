package tail

import (
	"errors"
	"testing"
)

// --- test helpers ---

type recorder struct {
	lines  []string
	closed bool
	fail   bool
}

func (r *recorder) Write(line string) error {
	if r.fail {
		return errors.New("write error")
	}
	r.lines = append(r.lines, line)
	return nil
}

func (r *recorder) Close() error {
	r.closed = true
	return nil
}

// --- tests ---

func TestLimiter_DisabledWhenZero(t *testing.T) {
	rec := &recorder{}
	l := New(rec, 0)
	for _, line := range []string{"a", "b", "c"} {
		if err := l.Write(line); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	l.Close()
	if len(rec.lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(rec.lines))
	}
}

func TestLimiter_DisabledWhenNegative(t *testing.T) {
	rec := &recorder{}
	l := New(rec, -5)
	l.Write("x")
	l.Close()
	if len(rec.lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(rec.lines))
	}
}

func TestLimiter_BuffersAndFlushesOnClose(t *testing.T) {
	rec := &recorder{}
	l := New(rec, 3)
	for _, line := range []string{"a", "b", "c"} {
		l.Write(line)
	}
	if len(rec.lines) != 0 {
		t.Fatal("expected no lines before Close")
	}
	l.Close()
	if len(rec.lines) != 3 {
		t.Fatalf("expected 3 lines after Close, got %d", len(rec.lines))
	}
	if rec.lines[0] != "a" || rec.lines[2] != "c" {
		t.Fatalf("unexpected order: %v", rec.lines)
	}
}

func TestLimiter_KeepsLastN(t *testing.T) {
	rec := &recorder{}
	l := New(rec, 3)
	for _, line := range []string{"1", "2", "3", "4", "5"} {
		l.Write(line)
	}
	l.Close()
	want := []string{"3", "4", "5"}
	if len(rec.lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(rec.lines), rec.lines)
	}
	for i, w := range want {
		if rec.lines[i] != w {
			t.Errorf("line[%d]: want %q got %q", i, w, rec.lines[i])
		}
	}
}

func TestLimiter_ExactlyMax(t *testing.T) {
	rec := &recorder{}
	l := New(rec, 4)
	for _, line := range []string{"p", "q", "r", "s"} {
		l.Write(line)
	}
	l.Close()
	if len(rec.lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(rec.lines))
	}
}

func TestLimiter_CloseMarksUnderlying(t *testing.T) {
	rec := &recorder{}
	l := New(rec, 2)
	l.Write("only")
	l.Close()
	if !rec.closed {
		t.Fatal("expected underlying writer to be closed")
	}
}
