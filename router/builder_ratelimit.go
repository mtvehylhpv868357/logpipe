package router

import (
	"time"

	"github.com/your-org/logpipe/config"
	"github.com/your-org/logpipe/ratelimit"
)

// limiterForOutput constructs a rate Limiter from the output config.
// Returns nil if no rate limit is configured.
func limiterForOutput(o config.Output) *ratelimit.Limiter {
	rl := o.RateLimit
	if rl == nil || rl.MaxLines <= 0 {
		return nil
	}

	window := time.Second
	if rl.Window != "" {
		d, err := time.ParseDuration(rl.Window)
		if err == nil && d > 0 {
			window = d
		}
	}

	return ratelimit.New(rl.MaxLines, window)
}
