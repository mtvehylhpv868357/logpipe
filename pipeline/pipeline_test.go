package pipeline_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/yourorg/logpipe/metrics"
	"github.com/yourorg/logpipe/pipeline"
	"github.com/yourorg/logpipe/router"
)

// stubRouter records routed lines and controls whether Route returns true.
type stubRouter struct {
	mu     sync.Mutex
	lines  []string
	accept bool
}

func (s *stubRouter) Route(line string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lines = append(s.lines, line)
	return s.accept
}

func newRouter(accept bool) (*router.Router, *stubRouter) {
	stub := &stubRouter{accept: accept}
	// Build a real router with no routes so every line is dropped,
	// but wrap via the stub for assertion purposes.
	_ = stub
	// For unit testing we rely on a minimal router.
	rt := router.New(nil)
	return rt, stub
}

func TestPipeline_Run_CountsReceived(t *testing.T) {
	input := "line1\nline2\nline3\n"
	r := strings.NewReader(input)
	rt := router.New(nil)
	c := metrics.NewCounters()

	p := pipeline.New(r, rt, c)
	if err := p.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := c.Snapshot()
	if snap.Received != 3 {
		t.Errorf("expected 3 received, got %d", snap.Received)
	}
}

func TestPipeline_Run_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	rt := router.New(nil)
	c := metrics.NewCounters()

	p := pipeline.New(r, rt, c)
	if err := p.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := c.Snapshot()
	if snap.Received != 0 {
		t.Errorf("expected 0 received, got %d", snap.Received)
	}
}

func TestPipeline_Run_DropsUnmatchedLines(t *testing.T) {
	input := "alpha\nbeta\n"
	r := strings.NewReader(input)
	// Router with no routes drops everything.
	rt := router.New(nil)
	c := metrics.NewCounters()

	p := pipeline.New(r, rt, c)
	if err := p.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := c.Snapshot()
	if snap.Received != 2 {
		t.Errorf("expected 2 received, got %d", snap.Received)
	}
	if snap.Dropped != 2 {
		t.Errorf("expected 2 dropped, got %d", snap.Dropped)
	}
}
