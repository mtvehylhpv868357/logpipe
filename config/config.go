package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TransformConfig describes a single transformation step.
type TransformConfig struct {
	Type    string            `yaml:"type"`
	Options map[string]string `yaml:"options"`
}

// OutputConfig describes a single output destination.
type OutputConfig struct {
	Type       string            `yaml:"type"`
	Path       string            `yaml:"path,omitempty"`
	Filters    []FilterConfig    `yaml:"filters,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
}

// FilterConfig describes a single filter rule.
type FilterConfig struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// Config is the top-level configuration structure.
type Config struct {
	Outputs []OutputConfig `yaml:"outputs"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: invalid YAML: %w", err)
	}

	if len(cfg.Outputs) == 0 {
		return nil, fmt.Errorf("config: at least one output is required")
	}

	for i, o := range cfg.Outputs {
		if o.Type == "" {
			return nil, fmt.Errorf("config: output[%d] missing type", i)
		}
	}

	return &cfg, nil
}
