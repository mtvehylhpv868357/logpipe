package multiline

import (
	"testing"
	"time"
)

func TestNew_InvalidStartPattern(t *testing.T) {
	_, err := New("[", "")
	if err == nil {
		t.Fatal("expected error for invalid start pattern")
	}
}

func TestNew_InvalidContinuePattern(t *testing.T) {
	_, err := New("^ERROR", "[")
	if err == nil {
		t.Fatal("expected error for invalid continue pattern")
	}
}

func TestAggregator_SingleLineEvents(t *testing.T) {
	a, _ := New("^ERROR", "")
	if event, ok := a.Add("ERROR foo"); ok {
		t.Fatalf("unexpected flush on first line: %q", event)
	}
	event, ok := a.Add("ERROR bar")
	if !ok {
		t.Fatal("expected flush when new event starts")
	}
	if event != "ERROR foo" {
		t.Fatalf("got %q want %q", event, "ERROR foo")
	}
}

func TestAggregator_MultiLineStack(t *testing.T) {
	a, _ := New("^ERROR", `^\s+at `)
	a.Add("ERROR something went wrong")
	a.Add("\tat foo.go:10")
	a.Add("\tat bar.go:20")

	event, ok := a.Add("ERROR next error")
	if !ok {
		t.Fatal("expected flush")
	}
	want := "ERROR something went wrong\n\tat foo.go:10\n\tat bar.go:20"
	if event != want {
		t.Fatalf("got %q want %q", event, want)
	}
}

func TestAggregator_Flush(t *testing.T) {
	a, _ := New("^ERROR", "")
	a.Add("ERROR orphan")
	event, ok := a.Flush()
	if !ok {
		t.Fatal("expected flush")
	}
	if event != "ERROR orphan" {
		t.Fatalf("got %q", event)
	}
	_, ok = a.Flush()
	if ok {
		t.Fatal("second flush should be empty")
	}
}

func TestAggregator_Timeout(t *testing.T) {
	now := time.Unix(1000, 0)
	a, _ := New("^ERROR", "", WithTimeout(2*time.Second))
	a.clock = func() time.Time { return now }

	a.Add("ERROR first")

	// Advance clock past timeout.
	now = now.Add(3 * time.Second)
	event, ok := a.Add("ERROR second")
	if !ok {
		t.Fatal("expected timeout flush")
	}
	if event != "ERROR first" {
		t.Fatalf("got %q", event)
	}
}

func TestAggregator_WithJoin(t *testing.T) {
	a, _ := New("^LOG", "", WithJoin(" | "))
	a.Add("LOG a")
	a.Add("LOG b")
	event, _ := a.Flush()
	// Only one line in buffer since second LOG triggered a flush of the first.
	if event != "LOG b" {
		t.Fatalf("got %q", event)
	}
}
