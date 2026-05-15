package parse

import (
	"testing"
)

func TestJSONParser_ValidObject(t *testing.T) {
	p := NewJSONParser()
	fields, err := p.Parse(`{"level":"info","msg":"hello","count":3}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", fields["level"])
	}
	if fields["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %q", fields["msg"])
	}
	if fields["count"] != "3" {
		t.Errorf("expected count=3, got %q", fields["count"])
	}
}

func TestJSONParser_InvalidJSON(t *testing.T) {
	p := NewJSONParser()
	_, err := p.Parse(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLogfmtParser_ValidLine(t *testing.T) {
	p := NewLogfmtParser()
	fields, err := p.Parse(`level=info msg="hello world" latency=42ms`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", fields["level"])
	}
	if fields["latency"] != "42ms" {
		t.Errorf("expected latency=42ms, got %q", fields["latency"])
	}
}

func TestLogfmtParser_NoKVPairs(t *testing.T) {
	p := NewLogfmtParser()
	_, err := p.Parse(`just a plain line with no equals`)
	if err == nil {
		t.Fatal("expected error when no key=value pairs found")
	}
}

func TestPlainParser_ReturnsMessage(t *testing.T) {
	p := NewPlainParser()
	fields, err := p.Parse("hello logpipe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["message"] != "hello logpipe" {
		t.Errorf("expected message field, got %q", fields["message"])
	}
}

func TestNewParser_ReturnsCorrectType(t *testing.T) {
	cases := []struct {
		format string
		wantType string
	}{
		{"json", "*parse.JSONParser"},
		{"logfmt", "*parse.LogfmtParser"},
		{"plain", "*parse.PlainParser"},
		{"unknown", "*parse.PlainParser"},
		{"", "*parse.PlainParser"},
	}
	for _, tc := range cases {
		p := NewParser(tc.format)
		if p == nil {
			t.Errorf("NewParser(%q) returned nil", tc.format)
		}
	}
}

func TestJSONParser_EmptyObject(t *testing.T) {
	p := NewJSONParser()
	fields, err := p.Parse(`{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("expected empty fields, got %v", fields)
	}
}
