// Package lineprefix prepends a static or dynamic prefix string to each line.
package lineprefix

import "fmt"

// Prepender prepends a fixed prefix to every line it processes.
type Prepender struct {
	prefix string
}

// New returns a Prepender that prepends prefix to each line.
// If prefix is empty the Prepender is a no-op.
func New(prefix string) *Prepender {
	return &Prepender{prefix: prefix}
}

// Transform prepends the configured prefix to line.
// If the prefix is empty the original line is returned unchanged.
func (p *Prepender) Transform(line string) string {
	if p.prefix == "" {
		return line
	}
	return fmt.Sprintf("%s%s", p.prefix, line)
}

// DynamicPrepender prepends a prefix produced by a user-supplied function.
type DynamicPrepender struct {
	fn func() string
}

// NewDynamic returns a DynamicPrepender that calls fn() before each line to
// obtain the prefix to prepend. fn must not be nil.
func NewDynamic(fn func() string) *DynamicPrepender {
	if fn == nil {
		panic("lineprefix: fn must not be nil")
	}
	return &DynamicPrepender{fn: fn}
}

// Transform calls the stored function to obtain the current prefix and
// prepends it to line.
func (d *DynamicPrepender) Transform(line string) string {
	prefix := d.fn()
	if prefix == "" {
		return line
	}
	return fmt.Sprintf("%s%s", prefix, line)
}
