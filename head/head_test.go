package head_test

import (
	"testing"

	"github.com/yourorg/logpipe/head"
)

func TestLimiter_DisabledWhenZero(t *testing.T) {
	l := head.New(0)
	for i := 0; i < 100; i++ {
		if !l.Allow() {
			t.Fatal("expected all lines to pass when max=0")
		}
	}
}

func TestLimiter_DisabledWhenNegative(t *testing.T) {
	l := head.New(-5)
	for i := 0; i < 10; i++ {
		if !l.Allow() {
			t.Fatal("expected all lines to pass when max<0")
		}
	}
}

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l := head.New(3)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected line %d to be allowed", i+1)
		}
	}
}

func TestLimiter_DropsAfterMax(t *testing.T) {
	l := head.New(3)
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	for i := 0; i < 5; i++ {
		if l.Allow() {
			t.Fatal("expected lines beyond max to be dropped")
		}
	}
}

func TestLimiter_ResetRestoresAllowance(t *testing.T) {
	l := head.New(2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected line to be dropped before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected line to be allowed after reset")
	}
}

func TestLimiter_Remaining_Disabled(t *testing.T) {
	l := head.New(0)
	if got := l.Remaining(); got != -1 {
		t.Fatalf("expected -1 for disabled limiter, got %d", got)
	}
}

func TestLimiter_Remaining_DecreasesWithAllows(t *testing.T) {
	l := head.New(5)
	if r := l.Remaining(); r != 5 {
		t.Fatalf("expected 5, got %d", r)
	}
	l.Allow()
	l.Allow()
	if r := l.Remaining(); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestLimiter_Remaining_NeverBelowZero(t *testing.T) {
	l := head.New(1)
	l.Allow()
	l.Allow()
	l.Allow()
	if r := l.Remaining(); r != 0 {
		t.Fatalf("expected 0, got %d", r)
	}
}
