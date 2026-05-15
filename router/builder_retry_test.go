package router

import (
	"errors"
	"testing"

	"github.com/yourorg/logpipe/config"
)

type stubRetryWriter struct {
	calls int
	fail  bool
}

func (s *stubRetryWriter) Write(line string) error {
	s.calls++
	if s.fail {
		return errors.New("fail")
	}
	return nil
}

func (s *stubRetryWriter) Close() error { return nil }

func TestRetryWriterForOutput_NilWhenNoRetry(t *testing.T) {
	inner := &stubRetryWriter{}
	out := config.Output{Type: "stdout"}
	w := retryWriterForOutput(inner, out)
	if w != inner {
		t.Error("expected same writer when no retry config")
	}
}

func TestRetryWriterForOutput_WrapsWhenConfigPresent(t *testing.T) {
	inner := &stubRetryWriter{fail: true}
	out := config.Output{
		Type: "stdout",
		Retry: &config.RetryConfig{MaxAttempts: 3, DelayMs: 0},
	}
	w := retryWriterForOutput(inner, out)
	if w == inner {
		t.Error("expected wrapped writer")
	}
	// All 3 attempts should fail; calls == 3
	err := w.Write("test")
	if err == nil {
		t.Fatal("expected error from all retries failing")
	}
	if inner.calls != 3 {
		t.Errorf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryWriterForOutput_DefaultsMinAttempts(t *testing.T) {
	inner := &stubRetryWriter{fail: true}
	out := config.Output{
		Type:  "stdout",
		Retry: &config.RetryConfig{MaxAttempts: 0},
	}
	w := retryWriterForOutput(inner, out)
	w.Write("x")
	if inner.calls != 1 {
		t.Errorf("expected 1 call with zero attempts, got %d", inner.calls)
	}
}
