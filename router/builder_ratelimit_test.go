package router

import (
	"testing"

	"github.com/your-org/logpipe/config"
)

func rateLimitOutput(maxLines int, window string) config.Output {
	return config.Output{
		Type: "stdout",
		RateLimit: &config.RateLimitConfig{
			MaxLines: maxLines,
			Window:   window,
		},
	}
}

func TestLimiterForOutput_Nil_WhenNoRateLimit(t *testing.T) {
	o := config.Output{Type: "stdout"}
	if limiterForOutput(o) != nil {
		t.Fatal("expected nil limiter when RateLimit is not set")
	}
}

func TestLimiterForOutput_Nil_WhenMaxLinesZero(t *testing.T) {
	o := rateLimitOutput(0, "1s")
	if limiterForOutput(o) != nil {
		t.Fatal("expected nil limiter when MaxLines is 0")
	}
}

func TestLimiterForOutput_DefaultsToOneSecondWindow(t *testing.T) {
	o := rateLimitOutput(10, "")
	l := limiterForOutput(o)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	// Allow up to 10 calls, 11th should be blocked
	for i := 0; i < 10; i++ {
		if !l.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected block after max")
	}
}

func TestLimiterForOutput_InvalidWindowFallsBackToSecond(t *testing.T) {
	o := rateLimitOutput(5, "not-a-duration")
	l := limiterForOutput(o)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected block after max with fallback window")
	}
}

func TestLimiterForOutput_CustomWindow(t *testing.T) {
	o := rateLimitOutput(3, "500ms")
	l := limiterForOutput(o)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected block after max")
	}
}
