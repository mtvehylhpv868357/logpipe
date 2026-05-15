// Package truncate provides a transformer that truncates log lines exceeding a maximum byte length.
package truncate

import "fmt"

const defaultSuffix = "...[truncated]"

// Truncator truncates lines that exceed MaxLen bytes.
type Truncator struct {
	MaxLen int
	Suffix string
}

// New creates a Truncator with the given max length and default suffix.
// If maxLen <= 0, the Truncator is disabled and lines pass through unchanged.
func New(maxLen int) *Truncator {
	return &Truncator{
		MaxLen: maxLen,
		Suffix: defaultSuffix,
	}
}

// NewWithSuffix creates a Truncator with a custom suffix appended to truncated lines.
func NewWithSuffix(maxLen int, suffix string) *Truncator {
	return &Truncator{
		MaxLen: maxLen,
		Suffix: suffix,
	}
}

// Apply truncates the line if it exceeds MaxLen bytes.
// If MaxLen <= 0 the line is returned unchanged.
func (t *Truncator) Apply(line string) string {
	if t.MaxLen <= 0 {
		return line
	}
	if len(line) <= t.MaxLen {
		return line
	}
	cutAt := t.MaxLen - len(t.Suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return fmt.Sprintf("%s%s", line[:cutAt], t.Suffix)
}
