package retry

import (
	"errors"
	"testing"
	"time"
)

type mockWriter struct {
	calls    int
	failUntil int
	closed  bool
	written []string
}

func (m *mockWriter) Write(line string) error {
	m.calls++
	if m.calls <= m.failUntil {
		return errors.New("write error")
	}
	m.written = append(m.written, line)
	return nil
}

func (m *mockWriter) Close() error {
	m.closed = true
	return nil
}

func noSleep(_ time.Duration) {}

func TestRetryWriter_SucceedsFirstAttempt(t *testing.T) {
	w := &mockWriter{}
	r := newWithSleep(w, 3, 0, noSleep)
	if err := r.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.calls != 1 {
		t.Errorf("expected 1 call, got %d", w.calls)
	}
}

func TestRetryWriter_RetriesOnFailure(t *testing.T) {
	w := &mockWriter{failUntil: 2}
	r := newWithSleep(w, 5, 0, noSleep)
	if err := r.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.calls != 3 {
		t.Errorf("expected 3 calls, got %d", w.calls)
	}
}

func TestRetryWriter_ReturnsLastError(t *testing.T) {
	w := &mockWriter{failUntil: 10}
	r := newWithSleep(w, 3, 0, noSleep)
	err := r.Write("hello")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if w.calls != 3 {
		t.Errorf("expected 3 calls, got %d", w.calls)
	}
}

func TestRetryWriter_MinAttemptsIsOne(t *testing.T) {
	w := &mockWriter{failUntil: 10}
	r := newWithSleep(w, 0, 0, noSleep)
	err := r.Write("x")
	if err == nil {
		t.Fatal("expected error")
	}
	if w.calls != 1 {
		t.Errorf("expected 1 call, got %d", w.calls)
	}
}

func TestRetryWriter_SleepsBetweenAttempts(t *testing.T) {
	w := &mockWriter{failUntil: 2}
	var slept []time.Duration
	sleepFn := func(d time.Duration) { slept = append(slept, d) }
	r := newWithSleep(w, 3, 50*time.Millisecond, sleepFn)
	r.Write("line")
	if len(slept) != 2 {
		t.Errorf("expected 2 sleeps, got %d", len(slept))
	}
}

func TestRetryWriter_Close(t *testing.T) {
	w := &mockWriter{}
	r := New(w, 1, 0)
	r.Close()
	if !w.closed {
		t.Error("expected inner writer to be closed")
	}
}
