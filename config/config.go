package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// OutputConfig describes a single output destination.
type OutputConfig struct {
	Type      string            `yaml:"type"`
	Path      string            `yaml:"path,omitempty"`
	Filters   []FilterConfig    `yaml:"filters,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Sampling  *SamplingConfig   `yaml:"sampling,omitempty"`
	RateLimit *RateLimitConfig  `yaml:"rate_limit,omitempty"`
	Buffer    *BufferConfig     `yaml:"buffer,omitempty"`
}

// BufferConfig holds batching parameters for an output.
type BufferConfig struct {
	MaxSize  int    `yaml:"max_size"`
	Interval string `yaml:"interval"`
}

// Interval parses the interval string, defaulting to 1s on error.
func (b *BufferConfig) Interval() time.Duration {
	if b.Interval == "" {
		return 0
	}
	d, err := time.ParseDuration(b.Interval)
	if err != nil {
		return time.Second
	}
	return d
}

type FilterConfig struct {
	Type    string `yaml:"type"`
	Match   string `yaml:"match,omitempty"`
	Pattern string `yaml:"pattern,omitempty"`
	Level   string `yaml:"level,omitempty"`
}

type TransformConfig struct {
	Type   string `yaml:"type"`
	Prefix string `yaml:"prefix,omitempty"`
	Format string `yaml:"format,omitempty"`
}

type SamplingConfig struct {
	Type string  `yaml:"type"`
	Rate float64 `yaml:"rate,omitempty"`
	N    int     `yaml:"n,omitempty"`
}

type RateLimitConfig struct {
	MaxLines int    `yaml:"max_lines"`
	Window   string `yaml:"window,omitempty"`
}

// Config is the top-level configuration.
type Config struct {
	Outputs []OutputConfig `yaml:"outputs"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Outputs) == 0 {
		return nil, errors.New("config: no outputs defined")
	}
	return &cfg, nil
}
