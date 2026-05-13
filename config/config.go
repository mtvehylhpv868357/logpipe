package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Output defines a single output destination with optional filtering rules.
type Output struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`    // "stdout", "file", "http"
	Target  string   `yaml:"target"` // file path or URL depending on type
	Filters []Filter `yaml:"filters"`
}

// Filter defines a rule to match log lines.
type Filter struct {
	Field   string `yaml:"field"`   // "level", "message", "any"
	Match   string `yaml:"match"`   // substring or regex pattern
	Negate  bool   `yaml:"negate"`  // invert the match
}

// Config is the top-level configuration for logpipe.
type Config struct {
	Outputs []Output `yaml:"outputs"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks that the config contains at least one valid output.
func (c *Config) validate() error {
	if len(c.Outputs) == 0 {
		return fmt.Errorf("at least one output must be defined")
	}

	validTypes := map[string]bool{"stdout": true, "file": true, "http": true}
	for i, o := range c.Outputs {
		if o.Name == "" {
			return fmt.Errorf("output[%d]: name is required", i)
		}
		if !validTypes[o.Type] {
			return fmt.Errorf("output[%d] %q: unsupported type %q", i, o.Name, o.Type)
		}
		if o.Type != "stdout" && o.Target == "" {
			return fmt.Errorf("output[%d] %q: target is required for type %q", i, o.Name, o.Type)
		}
	}

	return nil
}
