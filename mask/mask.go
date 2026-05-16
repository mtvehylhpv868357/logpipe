// Package mask provides field-level masking for structured log lines.
// It replaces the value of named fields with a fixed mask string.
package mask

import (
	"strings"
)

const defaultMask = "***"

// FieldMasker replaces values of specific keys in key=value formatted lines.
type FieldMasker struct {
	fields map[string]struct{}
	mask   string
}

// New returns a FieldMasker that replaces values for the given field names.
func New(fields []string) *FieldMasker {
	return NewWithMask(fields, defaultMask)
}

// NewWithMask returns a FieldMasker using a custom mask string.
func NewWithMask(fields []string, mask string) *FieldMasker {
	if mask == "" {
		mask = defaultMask
	}
	fm := &FieldMasker{
		fields: make(map[string]struct{}, len(fields)),
		mask:   mask,
	}
	for _, f := range fields {
		fm.fields[strings.ToLower(f)] = struct{}{}
	}
	return fm
}

// Apply scans the line for key=value or key="value" tokens and masks
// the values of any matching field names. Unrecognised tokens are passed through.
func (fm *FieldMasker) Apply(line string) string {
	if len(fm.fields) == 0 {
		return line
	}
	tokens := strings.Fields(line)
	for i, tok := range tokens {
		eqIdx := strings.IndexByte(tok, '=')
		if eqIdx < 1 {
			continue
		}
		key := strings.ToLower(tok[:eqIdx])
		if _, ok := fm.fields[key]; !ok {
			continue
		}
		val := tok[eqIdx+1:]
		if strings.HasPrefix(val, "\"") {
			tokens[i] = tok[:eqIdx+1] + "\"" + fm.mask + "\""
		} else {
			tokens[i] = tok[:eqIdx+1] + fm.mask
		}
	}
	return strings.Join(tokens, " ")
}
