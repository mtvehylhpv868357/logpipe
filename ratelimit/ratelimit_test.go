package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l := New(3, time.Second)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() false after max reached")
	}
}

func TestLimiter_ResetsAfterWindow(t *testing.T) {
	now := time.Unix(1000, 0)
	clock := func() time.Time { return now }
	l := newWithClock(2, 100*time.Millisecond, clock)

	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected block after max")
	}

	now = now.Add(200 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected Allow() true after window reset")
	}
}

func TestLimiter_DisabledWhenMaxZero(t *testing.T) {
	l := New(0, time.Second)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatal("expected all lines to pass when max=0")
		}
	}
}

func TestLimiter_DisabledWhenMaxNegative(t *testing.T) {
	l := New(-5, time.Second)
	if !l.Allow() {
		t.Fatal("expected Allow() true for negative max")
	}
}

func TestLimiter_Reset(t *testing.T) {
	l := New(1, time.Hour)
	l.Allow()
	if l.Allow() {
		t.Fatal("expected block before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected Allow() true after Reset")
	}
}

func TestLimiter_ConcurrentSafe(t *testing.T) {
	l := New(500, time.Second)
	var wg sync.WaitGroup
	allowed := make(chan bool, 1000)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- l.Allow()
		}()
	}
	wg.Wait()
	close(allowed)
	count := 0
	for a := range allowed {
		if a {
			count++
		}
	}
	if count != 500 {
		t.Fatalf("expected exactly 500 allowed, got %d", count)
	}
}
