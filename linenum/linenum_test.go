package linenum

import (
	"fmt"
	"sync"
	"testing"
)

func TestNumberer_DefaultFormat(t *testing.T) {
	n := New()
	got := n.Transform("hello")
	want := "[1] hello"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNumberer_IncrementsPerCall(t *testing.T) {
	n := New()
	for i := 1; i <= 5; i++ {
		got := n.Transform("line")
		want := fmt.Sprintf("[%d] line", i)
		if got != want {
			t.Errorf("call %d: got %q, want %q", i, got, want)
		}
	}
}

func TestNumberer_CustomFormat(t *testing.T) {
	n := NewWithFormat("%04d: ")
	got := n.Transform("msg")
	want := "0001: msg"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNumberer_EmptyFormatFallsBackToDefault(t *testing.T) {
	n := NewWithFormat("")
	got := n.Transform("x")
	want := "[1] x"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNumberer_Reset(t *testing.T) {
	n := New()
	n.Transform("a")
	n.Transform("b")
	n.Reset()
	got := n.Transform("c")
	want := "[1] c"
	if got != want {
		t.Errorf("after reset got %q, want %q", got, want)
	}
}

func TestNumberer_ResetRestoresCount(t *testing.T) {
	n := New()
	n.Transform("a")
	n.Transform("b")
	n.Reset()
	if n.Count() != 0 {
		t.Errorf("expected count 0 after reset, got %d", n.Count())
	}
}

func TestNumberer_Count(t *testing.T) {
	n := New()
	if n.Count() != 0 {
		t.Fatalf("expected initial count 0, got %d", n.Count())
	}
	n.Transform("a")
	n.Transform("b")
	if n.Count() != 2 {
		t.Fatalf("expected count 2, got %d", n.Count())
	}
}

func TestNumberer_ConcurrentSafety(t *testing.T) {
	n := New()
	var wg sync.WaitGroup
	const goroutines = 50
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			n.Transform("concurrent")
		}()
	}
	wg.Wait()
	if n.Count() != goroutines {
		t.Errorf("expected count %d, got %d", goroutines, n.Count())
	}
}
