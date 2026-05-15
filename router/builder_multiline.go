package router

import (
	"time"

	"github.com/user/logpipe/config"
	"github.com/user/logpipe/multiline"
)

// multilineAggForOutput constructs a multiline.Aggregator from the output
// config, or returns nil when multiline aggregation is not configured.
func multilineAggForOutput(o config.Output) *multiline.Aggregator {
	ml := o.Multiline
	if ml == nil || ml.StartPattern == "" {
		return nil
	}

	opts := []multiline.Option{}

	if ml.Join != "" {
		opts = append(opts, multiline.WithJoin(ml.Join))
	}

	if ml.TimeoutSeconds > 0 {
		opts = append(opts, multiline.WithTimeout(
			time.Duration(ml.TimeoutSeconds)*time.Second,
		))
	}

	agg, err := multiline.New(ml.StartPattern, ml.ContinuePattern, opts...)
	if err != nil {
		// Invalid patterns are a configuration error; log and disable.
		return nil
	}
	return agg
}
