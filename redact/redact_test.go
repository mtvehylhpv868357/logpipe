package redact

import (
	"testing"
)

func TestRegexRedactor_MasksMatch(t *testing.T) {
	r, err := NewRegexRedactor(`\d{4}-\d{4}-\d{4}-\d{4}`, "[CARD]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "payment card=1234-5678-9012-3456 processed"
	want := "payment card=[CARD] processed"
	if got := r.Redact(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRegexRedactor_DefaultMask(t *testing.T) {
	r, err := NewRegexRedactor(`secret`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Redact("my secret value")
	want := "my [REDACTED] value"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRegexRedactor_InvalidPattern(t *testing.T) {
	_, err := NewRegexRedactor(`[invalid`, "")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestRegexRedactor_NoMatch(t *testing.T) {
	r, err := NewRegexRedactor(`password=\S+`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "user logged in"
	if got := r.Redact(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestKeyValueRedactor_MasksValue(t *testing.T) {
	r := NewKeyValueRedactor("token", "")
	input := "auth token=abc123xyz request ok"
	got := r.Redact(input)
	want := "auth token=[REDACTED] request ok"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestKeyValueRedactor_CustomMask(t *testing.T) {
	r := NewKeyValueRedactor("password", "***")
	input := "login password=hunter2 failed"
	got := r.Redact(input)
	want := "login password=*** failed"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestKeyValueRedactor_NoMatch(t *testing.T) {
	r := NewKeyValueRedactor("secret", "")
	input := "nothing sensitive here"
	if got := r.Redact(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestChain_AppliesAllRedactors(t *testing.T) {
	r1, _ := NewRegexRedactor(`\d{3}-\d{2}-\d{4}`, "[SSN]")
	r2 := NewKeyValueRedactor("api_key", "")
	chain := NewChain(r1, r2)

	input := "ssn=123-45-6789 api_key=supersecret end"
	got := chain.Redact(input)

	if got == input {
		t.Error("expected line to be redacted")
	}
	if contains(got, "123-45-6789") {
		t.Error("SSN should be redacted")
	}
	if contains(got, "supersecret") {
		t.Error("api_key value should be redacted")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
