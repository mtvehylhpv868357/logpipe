package router

import (
	"fmt"

	"github.com/yourorg/logpipe/config"
	"github.com/yourorg/logpipe/filter"
	"github.com/yourorg/logpipe/output"
	"github.com/yourorg/logpipe/transform"
)

// Build constructs a Router from the provided Config.
func Build(cfg *config.Config) (*Router, error) {
	var routes []Route

	for i, oc := range cfg.Outputs {
		w, err := writerForOutput(oc)
		if err != nil {
			return nil, fmt.Errorf("builder: output[%d]: %w", i, err)
		}

		chain, err := filterChainForOutput(oc)
		if err != nil {
			return nil, fmt.Errorf("builder: output[%d] filters: %w", i, err)
		}

		tc, err := transformChainForOutput(oc)
		if err != nil {
			return nil, fmt.Errorf("builder: output[%d] transforms: %w", i, err)
		}

		routes = append(routes, Route{Filter: chain, Transform: tc, Writer: w})
	}

	return New(routes), nil
}

func writerForOutput(oc config.OutputConfig) (output.Writer, error) {
	switch oc.Type {
	case "stdout":
		return output.NewStdoutWriter(), nil
	case "file":
		if oc.Path == "" {
			return nil, fmt.Errorf("file output requires a path")
		}
		return output.NewFileWriter(oc.Path)
	default:
		return nil, fmt.Errorf("unknown output type %q", oc.Type)
	}
}

func filterChainForOutput(oc config.OutputConfig) (*filter.Chain, error) {
	var rules []filter.Rule
	for _, fc := range oc.Filters {
		r, err := filter.NewRule(fc.Type, fc.Value)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return filter.NewChain(rules...), nil
}

func transformChainForOutput(oc config.OutputConfig) (*transform.Chain, error) {
	var steps []transform.Transformer
	for _, tc := range oc.Transforms {
		t, err := transform.New(tc.Type, tc.Options)
		if err != nil {
			return nil, err
		}
		steps = append(steps, t)
	}
	return transform.NewChain(steps...), nil
}
