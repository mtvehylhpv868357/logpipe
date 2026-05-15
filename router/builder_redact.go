package router

import (
	"fmt"

	"github.com/aurc/logpipe/config"
	"github.com/aurc/logpipe/redact"
)

// redactChainForOutput builds a redact.Chain from the output's Redact config.
// Returns nil if no redactors are configured.
func redactChainForOutput(out config.Output) (*redact.Chain, error) {
	if len(out.Redact) == 0 {
		return nil, nil
	}

	var redactors []redact.Redactor

	for _, rc := range out.Redact {
		switch rc.Type {
		case "regex":
			if rc.Pattern == "" {
				return nil, fmt.Errorf("redact regex entry missing 'pattern'")
			}
			r, err := redact.NewRegexRedactor(rc.Pattern, rc.Mask)
			if err != nil {
				return nil, fmt.Errorf("invalid redact pattern %q: %w", rc.Pattern, err)
			}
			redactors = append(redactors, r)

		case "keyvalue", "kv":
			if rc.Key == "" {
				return nil, fmt.Errorf("redact keyvalue entry missing 'key'")
			}
			redactors = append(redactors, redact.NewKeyValueRedactor(rc.Key, rc.Mask))

		default:
			return nil, fmt.Errorf("unknown redact type %q (want 'regex' or 'keyvalue')", rc.Type)
		}
	}

	if len(redactors) == 0 {
		return nil, nil
	}
	return redact.NewChain(redactors...), nil
}
