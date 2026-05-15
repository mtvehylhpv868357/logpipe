package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// RetryConfig controls retry behaviour for a single output.
type RetryConfig struct {
	MaxAttempts int `yaml:"max_attempts"`
	DelayMs     int `yaml:"delay_ms"`
}

// Output describes a single log destination with its filters, transforms and
// optional features such as sampling, rate-limiting, buffering and retry.
type Output struct {
	Type       string            `yaml:"type"`
	Path       string            `yaml:"path,omitempty"`
	Filters    []map[string]string `yaml:"filters,omitempty"`
	Transforms []map[string]string `yaml:"transforms,omitempty"`
	Redact     []map[string]string `yaml:"redact,omitempty"`
	Sampling   *SamplingConfig    `yaml:"sampling,omitempty"`
	RateLimit  *RateLimitConfig   `yaml:"rate_limit,omitempty"`
	Buffer     *BufferConfig      `yaml:"buffer,omitempty"`
	Multiline  *MultilineConfig   `yaml:"multiline,omitempty"`
	Retry      *RetryConfig       `yaml:"retry,omitempty"`
}

// SamplingConfig selects a fraction of log lines.
type SamplingConfig struct {
	Rate float64 `yaml:"rate"`
	Nth  int     `yaml:"nth"`
}

// RateLimitConfig caps throughput to an output.
type RateLimitConfig struct {
	MaxLines int    `yaml:"max_lines"`
	Window   string `yaml:"window"`
}

// BufferConfig controls in-memory batching before writing.
type BufferConfig struct {
	MaxSize  int    `yaml:"max_size"`
	FlushMs  int    `yaml:"flush_ms"`
}

// MultilineConfig joins multi-line log events.
type MultilineConfig struct {
	StartPattern    string `yaml:"start_pattern"`
	ContinuePattern string `yaml:"continue_pattern"`
	Join            string `yaml:"join"`
	TimeoutMs       int    `yaml:"timeout_ms"`
}

// Config is the root configuration structure.
type Config struct {
	Outputs []Output `yaml:"outputs"`
}

// Load reads and parses a YAML config file at the given path.
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
