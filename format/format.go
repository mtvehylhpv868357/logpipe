// Package format provides log line formatters that serialize parsed fields
// into structured output formats such as JSON or logfmt.
package format

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Formatter transforms a map of parsed fields into a formatted string.
type Formatter interface {
	Format(fields map[string]string) string
}

// JSONFormatter serializes fields as a compact JSON object.
type JSONFormatter struct{}

// NewJSONFormatter returns a Formatter that outputs JSON.
func NewJSONFormatter() Formatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Format(fields map[string]string) string {
	if len(fields) == 0 {
		return "{}"
	}
	b, err := json.Marshal(fields)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// LogfmtFormatter serializes fields as key=value pairs sorted by key.
type LogfmtFormatter struct{}

// NewLogfmtFormatter returns a Formatter that outputs logfmt.
func NewLogfmtFormatter() Formatter {
	return &LogfmtFormatter{}
}

func (f *LogfmtFormatter) Format(fields map[string]string) string {
	if len(fields) == 0 {
		return ""
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := fields[k]
		if strings.ContainsAny(v, " \t\n\"") {
			v = fmt.Sprintf("%q", v)
		}
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, " ")
}

// TemplateFormatter formats fields using a simple {{key}} substitution template.
type TemplateFormatter struct {
	template string
}

// NewTemplateFormatter returns a Formatter that substitutes {{key}} placeholders.
func NewTemplateFormatter(tmpl string) Formatter {
	return &TemplateFormatter{template: tmpl}
}

func (f *TemplateFormatter) Format(fields map[string]string) string {
	out := f.template
	for k, v := range fields {
		out = strings.ReplaceAll(out, "{"+"{"+k+"}}", v)
	}
	return out
}
