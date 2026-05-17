package splitline_test

import (
	"testing"

	"github.com/yourorg/logpipe/splitline"
)

func TestSplitter_EmptyDelimiterPassesThrough(t *testing.T) {
	s := splitline.New("")
	got := s.Split("hello world")
	if len(got) != 1 || got[0] != "hello world" {
		t.Fatalf("expected original line, got %v", got)
	}
}

func TestSplitter_SingleSegment_NoDelimiterFound(t *testing.T) {
	s := splitline.New("|")
	got := s.Split("hello world")
	if len(got) != 1 || got[0] != "hello world" {
		t.Fatalf("expected original line, got %v", got)
	}
}

func TestSplitter_SplitsOnDelimiter(t *testing.T) {
	s := splitline.New("|")
	got := s.Split("a|b|c")
	if len(got) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(got), got)
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("unexpected parts: %v", got)
	}
}

func TestSplitter_TrimSpace(t *testing.T) {
	s := splitline.New(",", splitline.WithTrimSpace())
	got := s.Split(" foo , bar , baz ")
	if len(got) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(got), got)
	}
	if got[0] != "foo" || got[1] != "bar" || got[2] != "baz" {
		t.Errorf("unexpected parts: %v", got)
	}
}

func TestSplitter_EmptySegmentsDropped(t *testing.T) {
	s := splitline.New("|")
	got := s.Split("a||b")
	if len(got) != 2 {
		t.Fatalf("expected 2 parts after dropping empty, got %d: %v", len(got), got)
	}
	if got[0] != "a" || got[1] != "b" {
		t.Errorf("unexpected parts: %v", got)
	}
}

func TestSplitter_AllEmptySegmentsFallsBackToOriginal(t *testing.T) {
	s := splitline.New("|", splitline.WithTrimSpace())
	got := s.Split(" | | ")
	if len(got) != 1 || got[0] != " | | " {
		t.Fatalf("expected fallback to original line, got %v", got)
	}
}

func TestSplitter_MultiCharDelimiter(t *testing.T) {
	s := splitline.New("::")
	got := s.Split("key::value::extra")
	if len(got) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(got), got)
	}
	if got[0] != "key" || got[1] != "value" || got[2] != "extra" {
		t.Errorf("unexpected parts: %v", got)
	}
}
