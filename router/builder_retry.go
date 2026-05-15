package router

import (
	"time"

	"github.com/yourorg/logpipe/config"
	"github.com/yourorg/logpipe/retry"
)

// retryWriterForOutput wraps w in a RetryWriter if the output config specifies
// retry settings. Returns w unchanged if no retry config is present.
func retryWriterForOutput(w retry.Writer, out config.Output) retry.Writer {
	if out.Retry == nil {
		return w
	}

	attempts := out.Retry.MaxAttempts
	if attempts < 1 {
		attempts = 1
	}

	delay := time.Duration(out.Retry.DelayMs) * time.Millisecond
	return retry.New(w, attempts, delay)
}
