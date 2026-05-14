package metrics

import (
	"strings"
	"testing"
)

func TestTakeSnapshot_ReflectsCounters(t *testing.T) {
	c := NewCounters()
	c.IncReceived()
	c.IncReceived()
	c.IncRouted()
	c.IncDropped()
	c.IncErrors()

	s := c.TakeSnapshot()

	if s.Received != 2 {
		t.Errorf("expected Received=2, got %d", s.Received)
	}
	if s.Routed != 1 {
		t.Errorf("expected Routed=1, got %d", s.Routed)
	}
	if s.Dropped != 1 {
		t.Errorf("expected Dropped=1, got %d", s.Dropped)
	}
	if s.Errors != 1 {
		t.Errorf("expected Errors=1, got %d", s.Errors)
	}
}

func TestSnapshot_Print_ContainsFields(t *testing.T) {
	c := NewCounters()
	c.IncReceived()
	c.IncRouted()

	s := c.TakeSnapshot()
	var buf strings.Builder
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"received=1", "routed=1", "dropped=0", "errors=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}
