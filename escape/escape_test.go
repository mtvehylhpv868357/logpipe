package escape_test

import (
	"testing"

	"github.com/yourorg/logpipe/escape"
)

func TestSanitize_CleanLinePassesThrough(t *testing.T) {
	s := escape.New()
	input := "hello world"
	if got := s.Sanitize(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestSanitize_TabAndNewlinePreserved(t *testing.T) {
	s := escape.New()
	input := "col1\tcol2\nrow2"
	if got := s.Sanitize(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestSanitize_ControlCharReplaced(t *testing.T) {
	s := escape.New()
	// \x01 is a non-printable control character
	input := "hello\x01world"
	want := "hello\uFFFDworld"
	if got := s.Sanitize(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSanitize_CustomReplacement(t *testing.T) {
	s := escape.New(escape.WithReplacement('?'))
	input := "hi\x07there"
	want := "hi?there"
	if got := s.Sanitize(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSanitize_StripMode(t *testing.T) {
	s := escape.New(escape.WithStrip())
	input := "hi\x07\x08there"
	want := "hithere"
	if got := s.Sanitize(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSanitize_MultipleControlChars(t *testing.T) {
	s := escape.New(escape.WithReplacement('*'))
	input := "\x00\x01\x02"
	want := "***"
	if got := s.Sanitize(input); got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSanitize_EmptyString(t *testing.T) {
	s := escape.New()
	if got := s.Sanitize(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSanitize_UnicodePreserved(t *testing.T) {
	s := escape.New()
	input := "café \u4e2d\u6587"
	if got := s.Sanitize(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}
