// Package splitline provides a transformer that splits a single log line
// into multiple lines using a configurable delimiter.
package splitline

import "strings"

// Splitter splits a log line into multiple lines on a delimiter.
type Splitter struct {
	delimiter string
	trimSpace bool
}

// Option configures a Splitter.
type Option func(*Splitter)

// WithTrimSpace enables trimming of whitespace from each split segment.
func WithTrimSpace() Option {
	return func(s *Splitter) {
		s.trimSpace = true
	}
}

// New creates a Splitter that splits lines on the given delimiter.
// If delimiter is empty, the Splitter returns the original line unchanged.
func New(delimiter string, opts ...Option) *Splitter {
	s := &Splitter{delimiter: delimiter}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Split takes a single line and returns zero or more lines.
// If the delimiter is empty or not found, the original line is returned as-is.
func (s *Splitter) Split(line string) []string {
	if s.delimiter == "" {
		return []string{line}
	}

	parts := strings.Split(line, s.delimiter)

	out := parts[:0]
	for _, p := range parts {
		if s.trimSpace {
			p = strings.TrimSpace(p)
		}
		if p != "" {
			out = append(out, p)
		}
	}

	if len(out) == 0 {
		return []string{line}
	}
	return out
}
