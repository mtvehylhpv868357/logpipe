// Package escape provides utilities for sanitizing log lines by removing
// or replacing non-printable and control characters before forwarding.
package escape

import (
	"strings"
	"unicode"
)

// Sanitizer replaces or strips non-printable characters from log lines.
type Sanitizer struct {
	replacement rune
	strip       bool
}

// Option configures a Sanitizer.
type Option func(*Sanitizer)

// WithReplacement sets the rune used to replace non-printable characters.
// Defaults to the Unicode replacement character (U+FFFD).
func WithReplacement(r rune) Option {
	return func(s *Sanitizer) {
		s.replacement = r
		s.strip = false
	}
}

// WithStrip configures the Sanitizer to remove non-printable characters
// entirely instead of replacing them.
func WithStrip() Option {
	return func(s *Sanitizer) {
		s.strip = true
	}
}

// New creates a Sanitizer with the given options.
// By default it replaces non-printable characters with U+FFFD.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{
		replacement: unicode.ReplacementChar,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Sanitize processes a single log line, handling non-printable and control
// characters according to the Sanitizer configuration.
// Newline (\n) and tab (\t) characters are always preserved.
func (s *Sanitizer) Sanitize(line string) string {
	var b strings.Builder
	b.Grow(len(line))
	for _, r := range line {
		if r == '\n' || r == '\t' {
			b.WriteRune(r)
			continue
		}
		if unicode.IsControl(r) || !unicode.IsPrint(r) {
			if s.strip {
				continue
			}
			b.WriteRune(s.replacement)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
