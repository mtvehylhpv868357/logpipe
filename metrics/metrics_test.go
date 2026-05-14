package metrics

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestCounters_InitialValuesAreZero(t *testing.T) {
	var c Counters
	s := c.Snapshot()
	if s.LinesRead != 0 || s.LinesMatched != 0 || s.LinesDropped != 0 || s.WriteErrors != 0 {
		t.Errorf("expected all zero counters, got %+v", s)
	}
}

func TestCounters_Increment(t *testing.T) {
	var c Counters
	c.IncRead()
	c.IncRead()
	c.IncMatched()
	c.IncDropped()
	c.IncWriteError()

	s := c.Snapshot()
	if s.LinesRead != 2 {
		t.Errorf("LinesRead: want 2, got %d", s.LinesRead)
	}
	if s.LinesMatched != 1 {
		t.Errorf("LinesMatched: want 1, got %d", s.LinesMatched)
	}
	if s.LinesDropped != 1 {
		t.Errorf("LinesDropped: want 1, got %d", s.LinesDropped)
	}
	if s.WriteErrors != 1 {
		t.Errorf("WriteErrors: want 1, got %d", s.WriteErrors)
	}
}

func TestCounters_ConcurrentIncrements(t *testing.T) {
	var c Counters
	var wg sync.WaitGroup
	const goroutines = 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.IncRead()
			c.IncMatched()
		}()
	}
	wg.Wait()

	s := c.Snapshot()
	if s.LinesRead != goroutines {
		t.Errorf("LinesRead: want %d, got %d", goroutines, s.LinesRead)
	}
	if s.LinesMatched != goroutines {
		t.Errorf("LinesMatched: want %d, got %d", goroutines, s.LinesMatched)
	}
}

func TestSnapshot_Print(t *testing.T) {
	s := Snapshot{
		LinesRead:    10,
		LinesMatched: 7,
		LinesDropped: 3,
		WriteErrors:  1,
	}
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"10", "7", "3", "1", "lines read", "write errors"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q; got:\n%s", want, out)
		}
	}
}
