package transform

import (
	"strings"
	"testing"
)

func TestPrefixTransformer(t *testing.T) {
	p := &PrefixTransformer{Prefix: "[INFO] "}
	out := p.Apply("hello world")
	if out != "[INFO] hello world" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestUpperTransformer(t *testing.T) {
	u := &UpperTransformer{}
	out := u.Apply("hello world")
	if out != "HELLO WORLD" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTimestampTransformer_DefaultFormat(t *testing.T) {
	tr := &TimestampTransformer{}
	out := tr.Apply("msg")
	if !strings.Contains(out, "msg") {
		t.Errorf("output should contain original line, got: %q", out)
	}
	// RFC3339 timestamps contain 'T'
	if !strings.Contains(out, "T") {
		t.Errorf("output should contain RFC3339 timestamp, got: %q", out)
	}
}

func TestTimestampTransformer_CustomFormat(t *testing.T) {
	tr := &TimestampTransformer{Format: "2006"}
	out := tr.Apply("msg")
	parts := strings.SplitN(out, " ", 2)
	if len(parts) != 2 || len(parts[0]) != 4 {
		t.Errorf("expected 4-digit year prefix, got: %q", out)
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	c := NewChain(
		&PrefixTransformer{Prefix: ">> "},
		&UpperTransformer{},
	)
	out := c.Apply("hello")
	if out != ">> HELLO" {
		t.Errorf("unexpected chain output: %q", out)
	}
}

func TestChain_Empty(t *testing.T) {
	c := NewChain()
	out := c.Apply("unchanged")
	if out != "unchanged" {
		t.Errorf("empty chain should be identity, got: %q", out)
	}
}

func TestNew_KnownTypes(t *testing.T) {
	cases := []struct {
		kind string
		opts map[string]string
	}{
		{"prefix", map[string]string{"value": "X "}},
		{"timestamp", map[string]string{}},
		{"upper", map[string]string{}},
	}
	for _, tc := range cases {
		tr, err := New(tc.kind, tc.opts)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", tc.kind, err)
		}
		if tr == nil {
			t.Errorf("New(%q) returned nil", tc.kind)
		}
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("nonexistent", nil)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}
