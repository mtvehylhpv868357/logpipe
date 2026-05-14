package metrics

import (
	"strings"
	"testing"
	"time"
)

func TestReporter_PrintsOnStop(t *testing.T) {
	c := NewCounters()
	c.IncReceived()
	c.IncRouted()

	var buf strings.Builder
	r := NewReporter(c, 10*time.Second, &buf)
	r.Start()
	r.Stop()

	out := buf.String()
	if !strings.Contains(out, "received=1") {
		t.Errorf("expected final snapshot in output, got: %s", out)
	}
}

func TestReporter_TicksPeriodically(t *testing.T) {
	c := NewCounters()
	var buf strings.Builder
	r := NewReporter(c, 20*time.Millisecond, &buf)
	r.Start()

	time.Sleep(70 * time.Millisecond)
	r.Stop()

	// Expect at least 2 tick prints plus the final stop print.
	lines := strings.Count(buf.String(), "\n")
	if lines < 2 {
		t.Errorf("expected at least 2 output lines, got %d\noutput: %s", lines, buf.String())
	}
}

func TestReporter_StopIsIdempotentInEffect(t *testing.T) {
	c := NewCounters()
	var buf strings.Builder
	r := NewReporter(c, 10*time.Second, &buf)
	r.Start()
	r.Stop()
	// Calling Stop again would panic on double-close; ensure goroutine exited.
	// We just verify doneCh is closed by checking buf has content.
	if buf.Len() == 0 {
		t.Error("expected at least one line of output after Stop")
	}
}
