package transform

import (
	"fmt"
	"strings"
	"time"
)

// Transformer mutates a log line before forwarding.
type Transformer interface {
	Apply(line string) string
}

// PrefixTransformer prepends a static string to each line.
type PrefixTransformer struct {
	Prefix string
}

func (p *PrefixTransformer) Apply(line string) string {
	return p.Prefix + line
}

// TimestampTransformer prepends the current UTC time in RFC3339 format.
type TimestampTransformer struct {
	Format string
}

func (t *TimestampTransformer) Apply(line string) string {
	fmt := t.Format
	if fmt == "" {
		fmt = time.RFC3339
	}
	return time.Now().UTC().Format(fmt) + " " + line
}

// UpperTransformer converts each line to uppercase.
type UpperTransformer struct{}

func (u *UpperTransformer) Apply(line string) string {
	return strings.ToUpper(line)
}

// Chain applies a sequence of Transformers in order.
type Chain struct {
	steps []Transformer
}

// NewChain constructs a Chain from the given Transformers.
func NewChain(steps ...Transformer) *Chain {
	return &Chain{steps: steps}
}

// Apply runs the line through every step in sequence.
func (c *Chain) Apply(line string) string {
	for _, s := range c.steps {
		line = s.Apply(line)
	}
	return line
}

// New builds a Transformer from a type name and options map.
func New(kind string, opts map[string]string) (Transformer, error) {
	switch kind {
	case "prefix":
		return &PrefixTransformer{Prefix: opts["value"]}, nil
	case "timestamp":
		return &TimestampTransformer{Format: opts["format"]}, nil
	case "upper":
		return &UpperTransformer{}, nil
	default:
		return nil, fmt.Errorf("transform: unknown type %q", kind)
	}
}
