package truncate

import (
	"strings"
	"testing"
)

func TestTruncator_ShortLinePassesThrough(t *testing.T) {
	tr := New(50)
	input := "short line"
	got := tr.Apply(input)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestTruncator_ExactLengthPassesThrough(t *testing.T) {
	tr := New(10)
	input := "0123456789" // exactly 10 bytes
	got := tr.Apply(input)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestTruncator_LongLineIsTruncated(t *testing.T) {
	tr := New(20)
	input := strings.Repeat("a", 50)
	got := tr.Apply(input)
	if len(got) > 20 {
		t.Errorf("expected len <= 20, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, defaultSuffix) {
		t.Errorf("expected suffix %q in %q", defaultSuffix, got)
	}
}

func TestTruncator_CustomSuffix(t *testing.T) {
	tr := NewWithSuffix(15, ">>")
	input := strings.Repeat("b", 30)
	got := tr.Apply(input)
	if len(got) > 15 {
		t.Errorf("expected len <= 15, got %d", len(got))
	}
	if !strings.HasSuffix(got, ">>") {
		t.Errorf("expected custom suffix in %q", got)
	}
}

func TestTruncator_DisabledWhenMaxLenZero(t *testing.T) {
	tr := New(0)
	input := strings.Repeat("x", 200)
	got := tr.Apply(input)
	if got != input {
		t.Error("expected line to pass through unchanged when MaxLen is 0")
	}
}

func TestTruncator_DisabledWhenMaxLenNegative(t *testing.T) {
	tr := New(-5)
	input := strings.Repeat("y", 100)
	got := tr.Apply(input)
	if got != input {
		t.Error("expected line to pass through unchanged when MaxLen is negative")
	}
}

func TestTruncator_SuffixLongerThanMaxLen(t *testing.T) {
	// Suffix is longer than maxLen — cutAt clamps to 0, result is just the suffix.
	tr := NewWithSuffix(3, "...[truncated]")
	input := strings.Repeat("z", 50)
	got := tr.Apply(input)
	if got != "...[truncated]" {
		t.Errorf("expected only suffix, got %q", got)
	}
}
