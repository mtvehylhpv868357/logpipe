package timestamp_test

import (
	"errors"
	"testing"
	"time"

	"logpipe/timestamp"
)

// mockWriter captures written lines and can simulate errors.
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

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestStamper_PrependDefaultFormat(t *testing.T) {
	m := &mockWriter{}
	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	s := timestamp.NewWithFormat(m, "")
	// replace clock via internal helper exposed for tests
	_ = s // use New path below

	s2 := newStamperWithClock(m, "", fixedClock(fixed))
	if err := s2.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(m.lines))
	}
	want := "2024-06-01T12:00:00Z hello"
	if m.lines[0] != want {
		t.Errorf("got %q, want %q", m.lines[0], want)
	}
}

func TestStamper_CustomFormat(t *testing.T) {
	m := &mockWriter{}
	fixed := time.Date(2024, 1, 15, 8, 30, 0, 0, time.UTC)
	s := newStamperWithClock(m, "2006-01-02", fixedClock(fixed))
	_ = s.Write("msg")
	if m.lines[0] != "2024-01-15 msg" {
		t.Errorf("unexpected line: %q", m.lines[0])
	}
}

func TestStamper_WriteError(t *testing.T) {
	m := &mockWriter{fail: true}
	fixed := time.Now()
	s := newStamperWithClock(m, "", fixedClock(fixed))
	if err := s.Write("line"); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStamper_CloseDelegate(t *testing.T) {
	m := &mockWriter{}
	s := timestamp.New(m)
	if err := s.Close(); err != nil {
		t.Errorf("unexpected close error: %v", err)
	}
}

func TestStamper_MultipleWrites(t *testing.T) {
	m := &mockWriter{}
	fixed := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	s := newStamperWithClock(m, time.RFC3339, fixedClock(fixed))
	for _, line := range []string{"a", "b", "c"} {
		_ = s.Write(line)
	}
	if len(m.lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(m.lines))
	}
}

// newStamperWithClock is a test helper that wires a custom clock.
func newStamperWithClock(out interface {
	Write(string) error
	Close() error
}, format string, fn func() time.Time) *timestamp.Stamper {
	return timestamp.NewWithClockExported(out, format, fn)
}
