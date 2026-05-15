package retry

import (
	"time"
)

// Writer is an interface for something that can write a line and be closed.
type Writer interface {
	Write(line string) error
	Close() error
}

// RetryWriter wraps a Writer and retries failed writes up to MaxAttempts times.
type RetryWriter struct {
	inner    Writer
	max      int
	delay    time.Duration
	sleepFn  func(time.Duration)
}

// New returns a RetryWriter that retries up to maxAttempts times with the
// given delay between attempts.
func New(w Writer, maxAttempts int, delay time.Duration) *RetryWriter {
	return newWithSleep(w, maxAttempts, delay, time.Sleep)
}

func newWithSleep(w Writer, maxAttempts int, delay time.Duration, sleepFn func(time.Duration)) *RetryWriter {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &RetryWriter{
		inner:   w,
		max:     maxAttempts,
		delay:   delay,
		sleepFn: sleepFn,
	}
}

// Write attempts to write line, retrying on error up to the configured limit.
// Returns the last error if all attempts fail.
func (r *RetryWriter) Write(line string) error {
	var err error
	for i := 0; i < r.max; i++ {
		if err = r.inner.Write(line); err == nil {
			return nil
		}
		if i < r.max-1 && r.delay > 0 {
			r.sleepFn(r.delay)
		}
	}
	return err
}

// Close closes the underlying writer.
func (r *RetryWriter) Close() error {
	return r.inner.Close()
}
