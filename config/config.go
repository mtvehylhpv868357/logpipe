package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/yourorg/logpipe/filter"
)

// OutputType represents the destination type for log lines.
type OutputType string

const (
	OutputStdout OutputType = "stdout"
	OutputFile   OutputType = "file"
	OutputHTTP   OutputType = "http"
)

// Output defines a single output destination with optional filters.
type Output struct {
	Type    OutputType    `yaml:"type"`
	Target  string        `yaml:"target"`
	Filters []filter.Rule `yaml:"filters"`
}

// Config is the top-level configuration structure.
type Config struct {
	Outputs []Output `yaml:"outputs"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks the loaded config for semantic errors.
func (c *Config) validate() error {
	if len(c.Outputs) == 0 {
		return errors.New("config must define at least one output")
	}
	for i, out := range c.Outputs {
		switch out.Type {
		case OutputStdout, OutputFile, OutputHTTP:
			// valid
		default:
			return errors.New("output[" + string(rune('0'+i)) + "] has unknown type: " + string(out.Type))
		}
		if (out.Type == OutputFile || out.Type == OutputHTTP) && out.Target == "" {
			return errors.New("output target is required for type " + string(out.Type))
		}
		for j := range c.Outputs[i].Filters {
			if err := c.Outputs[i].Filters[j].Compile(); err != nil {
				return err
			}
		}
	}
	return nil
}
