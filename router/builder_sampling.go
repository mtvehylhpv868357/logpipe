package router

import (
	"fmt"
	"math/rand"

	"logpipe/config"
	"logpipe/sampling"
)

// samplerForOutput constructs a Sampler based on the sampling configuration
// embedded in an output definition. If no sampling config is present a
// PassthroughSampler is returned so all lines are forwarded.
func samplerForOutput(out config.Output) (sampling.Sampler, error) {
	sc := out.Sampling
	if sc == nil {
		return sampling.PassthroughSampler{}, nil
	}

	switch sc.Strategy {
	case "rate":
		if sc.Rate <= 0 || sc.Rate > 1.0 {
			return nil, fmt.Errorf(
				"output %q: sampling rate must be in (0, 1], got %v",
				out.Name, sc.Rate,
			)
		}
		src := rand.NewSource(sc.Seed)
		return sampling.NewRateSampler(sc.Rate, src), nil

	case "nth":
		if sc.N == 0 {
			return nil, fmt.Errorf(
				"output %q: sampling nth must be >= 1",
				out.Name,
			)
		}
		return sampling.NewNthSampler(sc.N), nil

	case "", "none":
		return sampling.PassthroughSampler{}, nil

	default:
		return nil, fmt.Errorf(
			"output %q: unknown sampling strategy %q",
			out.Name, sc.Strategy,
		)
	}
}
