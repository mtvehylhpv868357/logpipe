package router

import (
	"fmt"
	"logpipe/config"
	"logpipe/filter"
	"logpipe/output"
)

// Build constructs a Router from the application config.
// It creates the appropriate writer for each output and wires up filter chains.
func Build(cfg *config.Config) (*Router, error) {
	var routes []*Route

	for i, out := range cfg.Outputs {
		writer, err := writerForOutput(out)
		if err != nil {
			return nil, fmt.Errorf("output[%d]: %w", i, err)
		}

		chain := filter.NewChain(out.Rules)

		routes = append(routes, &Route{
			Chain:  chain,
			Writer: writer,
		})
	}

	return New(routes), nil
}

func writerForOutput(out config.Output) (interface {
	Write([]byte) (int, error)
	Close() error
}, error) {
	switch out.Type {
	case "stdout":
		return output.NewStdoutWriter(), nil
	case "file":
		if out.Path == "" {
			return nil, fmt.Errorf("file output requires a path")
		}
		w, err := output.NewFileWriter(out.Path)
		if err != nil {
			return nil, err
		}
		return w, nil
	default:
		return nil, fmt.Errorf("unknown output type: %q", out.Type)
	}
}
