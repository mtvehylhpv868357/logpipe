// Package drop provides a writer that discards lines matching a given pattern.
package drop

import (
	"fmt"
	"regexp"
)

// Dropper discards lines that match a compiled regular expression and forwards
// all other lines to the wrapped writer.
type Dropper struct {
	pattern *regexp.Regexp
	next    writer
}

type writer interface {
	Write(line string) error
	Close() error
}

// New compiles pattern and returns a Dropper that drops matching lines.
// Returns an error if pattern is not a valid regular expression.
func New(pattern string, next writer) (*Dropper, error) {
	if pattern == "" {
		return &Dropper{next: next}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("drop: invalid pattern %q: %w", pattern, err)
	}
	return &Dropper{pattern: re, next: next}, nil
}

// Write discards line if it matches the drop pattern; otherwise it forwards
// the line to the underlying writer.
func (d *Dropper) Write(line string) error {
	if d.pattern != nil && d.pattern.MatchString(line) {
		return nil
	}
	return d.next.Write(line)
}

// Close closes the underlying writer.
func (d *Dropper) Close() error {
	return d.next.Close()
}
