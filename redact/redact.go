// Package redact provides transformers that mask sensitive data in log lines.
package redact

import (
	"regexp"
	"strings"
)

// Redactor masks sensitive patterns in a log line.
type Redactor interface {
	Redact(line string) string
}

// RegexRedactor replaces all matches of a compiled regex with a mask string.
type RegexRedactor struct {
	pattern *regexp.Regexp
	mask    string
}

// NewRegexRedactor compiles pattern and returns a RegexRedactor.
// mask defaults to "[REDACTED]" when empty.
func NewRegexRedactor(pattern, mask string) (*RegexRedactor, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if mask == "" {
		mask = "[REDACTED]"
	}
	return &RegexRedactor{pattern: re, mask: mask}, nil
}

// Redact replaces all regex matches in line with the mask.
func (r *RegexRedactor) Redact(line string) string {
	return r.pattern.ReplaceAllString(line, r.mask)
}

// KeyValueRedactor masks the value portion of key=value pairs for a given key.
type KeyValueRedactor struct {
	key     string
	mask    string
	pattern *regexp.Regexp
}

// NewKeyValueRedactor creates a redactor that masks values for the given key.
// It matches patterns like key=value or key="value".
func NewKeyValueRedactor(key, mask string) *KeyValueRedactor {
	if mask == "" {
		mask = "[REDACTED]"
	}
	// matches key=value or key="..."
	pat := regexp.MustCompile(`(?i)(` + regexp.QuoteMeta(key) + `=)(["']?)[^\s"'&,;]*(\2)`)
	return &KeyValueRedactor{key: key, mask: mask, pattern: pat}
}

// Redact masks the value of the configured key in line.
func (k *KeyValueRedactor) Redact(line string) string {
	return k.pattern.ReplaceAllStringFunc(line, func(match string) string {
		eq := strings.Index(match, "=")
		if eq == -1 {
			return match
		}
		return match[:eq+1] + k.mask
	})
}

// Chain applies multiple Redactors in sequence.
type Chain struct {
	redactors []Redactor
}

// NewChain returns a Chain that applies each Redactor in order.
func NewChain(redactors ...Redactor) *Chain {
	return &Chain{redactors: redactors}
}

// Redact runs each redactor in the chain against line.
func (c *Chain) Redact(line string) string {
	for _, r := range c.redactors {
		line = r.Redact(line)
	}
	return line
}
