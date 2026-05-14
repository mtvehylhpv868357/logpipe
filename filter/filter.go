package filter

import (
	"regexp"
	"strings"
)

// Rule defines a single filtering rule applied to log lines.
type Rule struct {
	Contains string `yaml:"contains"`
	Regex    string `yaml:"regex"`
	Level    string `yaml:"level"`

	compiledRegex *regexp.Regexp
}

// Compile prepares the rule for matching (e.g., compiles regex).
func (r *Rule) Compile() error {
	if r.Regex != "" {
		re, err := regexp.Compile(r.Regex)
		if err != nil {
			return err
		}
		r.compiledRegex = re
	}
	return nil
}

// Match returns true if the log line satisfies the rule.
func (r *Rule) Match(line string) bool {
	if r.Contains != "" && !strings.Contains(line, r.Contains) {
		return false
	}
	if r.compiledRegex != nil && !r.compiledRegex.MatchString(line) {
		return false
	}
	if r.Level != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(r.Level)) {
		return false
	}
	return true
}

// Chain holds a list of rules and matches a line against all of them.
// All rules must match (AND semantics).
type Chain struct {
	Rules []Rule
}

// NewChain creates a Chain from a slice of rules, compiling each one.
func NewChain(rules []Rule) (*Chain, error) {
	for i := range rules {
		if err := rules[i].Compile(); err != nil {
			return nil, err
		}
	}
	return &Chain{Rules: rules}, nil
}

// Match returns true if the line matches every rule in the chain.
// An empty chain matches all lines.
func (c *Chain) Match(line string) bool {
	for _, r := range c.Rules {
		if !r.Match(line) {
			return false
		}
	}
	return true
}
