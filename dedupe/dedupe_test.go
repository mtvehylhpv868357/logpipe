package dedupe

import (
	"sync"
	"testing"
	"time"
)

func TestFilter_AllowsUniqueLines(t *testing.T) {
	f := New(time.Second)
	if !f.Allow("hello") {
		t.Fatal("expected first occurrence to be allowed")
	}
	if !f.Allow("world") {
		t.Fatal("expected different line to be allowed")
	}
}

func TestFilter_SuppressDuplicate(t *testing.T) {
	f := New(time.Second)
	f.Allow("dup")
	if f.Allow("dup") {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestFilter_AllowsAfterWindowExpires(t *testing.T) {
	now := time.Unix(1000, 0)
	clock := func() time.Time { return now }

	f := newWithClock(100*time.Millisecond, clock)
	f.Allow("line")

	// advance past the window
	now = now.Add(200 * time.Millisecond)
	if !f.Allow("line") {
		t.Fatal("expected line to be allowed after window expired")
	}
}

func TestFilter_DisabledWithZeroWindow(t *testing.T) {
	f := New(0)
	f.Allow("x")
	if !f.Allow("x") {
		t.Fatal("expected all lines to pass when window is zero")
	}
}

func TestFilter_DisabledWithNegativeWindow(t *testing.T) {
	f := New(-time.Second)
	f.Allow("x")
	if !f.Allow("x") {
		t.Fatal("expected all lines to pass when window is negative")
	}
}

func TestFilter_Reset(t *testing.T) {
	f := New(time.Minute)
	f.Allow("line")
	f.Reset()
	if !f.Allow("line") {
		t.Fatal("expected line to be allowed after reset")
	}
}

func TestFilter_ConcurrentAccess(t *testing.T) {
	f := New(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.Allow("concurrent-line")
		}()
	}
	wg.Wait()
}
