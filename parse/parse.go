// Package parse provides structured log line parsers (JSON, logfmt, plain text).
package parse

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fields represents key-value pairs extracted from a log line.
type Fields map[string]string

// Parser extracts structured fields from a raw log line.
type Parser interface {
	Parse(line string) (Fields, error)
}

// JSONParser parses JSON-formatted log lines.
type JSONParser struct{}

// NewJSONParser returns a new JSONParser.
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

// Parse decodes a JSON object from line and returns string fields.
func (p *JSONParser) Parse(line string) (Fields, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}
	fields := make(Fields, len(raw))
	for k, v := range raw {
		fields[k] = fmt.Sprintf("%v", v)
	}
	return fields, nil
}

// LogfmtParser parses logfmt-formatted log lines (key=value pairs).
type LogfmtParser struct{}

// NewLogfmtParser returns a new LogfmtParser.
func NewLogfmtParser() *LogfmtParser {
	return &LogfmtParser{}
}

// Parse decodes key=value pairs from a logfmt line.
func (p *LogfmtParser) Parse(line string) (Fields, error) {
	fields := make(Fields)
	for _, token := range strings.Fields(line) {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := parts[0]
		v := strings.Trim(parts[1], `"`)
		fields[k] = v
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("logfmt parse: no key=value pairs found")
	}
	return fields, nil
}

// PlainParser returns the whole line under the key "message".
type PlainParser struct{}

// NewPlainParser returns a new PlainParser.
func NewPlainParser() *PlainParser {
	return &PlainParser{}
}

// Parse wraps the raw line as a single "message" field.
func (p *PlainParser) Parse(line string) (Fields, error) {
	return Fields{"message": line}, nil
}

// NewParser returns a Parser for the given format name.
// Supported formats: "json", "logfmt", "plain" (default).
func NewParser(format string) Parser {
	switch strings.ToLower(format) {
	case "json":
		return NewJSONParser()
	case "logfmt":
		return NewLogfmtParser()
	default:
		return NewPlainParser()
	}
}
