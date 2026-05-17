// Package labelinject provides a transformer that injects static key=value
// labels into every log line passing through the pipeline.
package labelinject

import "fmt"

// Injector appends static label fields to each log line.
type Injector struct {
	labels []label
}

type label struct {
	key   string
	value string
}

// New returns an Injector that will append the given key/value pairs to every
// line. Keys and values must be non-empty; any pair where either is empty is
// silently ignored.
func New(pairs map[string]string) *Injector {
	inj := &Injector{}
	for k, v := range pairs {
		if k == "" || v == "" {
			continue
		}
		inj.labels = append(inj.labels, label{key: k, value: v})
	}
	return inj
}

// Transform appends each configured label as " key=value" to line.
// If no labels are configured the original line is returned unchanged.
func (inj *Injector) Transform(line string) string {
	if len(inj.labels) == 0 {
		return line
	}
	out := line
	for _, l := range inj.labels {
		out += fmt.Sprintf(" %s=%s", l.key, l.value)
	}
	return out
}
