package filter

import (
	"testing"
)

func TestRule_Contains(t *testing.T) {
	r := Rule{Contains: "ERROR"}
	if err := r.Compile(); err != nil {
		t.Fatal(err)
	}
	if !r.Match("2024/01/01 ERROR something broke") {
		t.Error("expected match for line containing ERROR")
	}
	if r.Match("2024/01/01 INFO all good") {
		t.Error("expected no match for line without ERROR")
	}
}

func TestRule_Regex(t *testing.T) {
	r := Rule{Regex: `\bWARN\b`}
	if err := r.Compile(); err != nil {
		t.Fatal(err)
	}
	if !r.Match("WARN: disk space low") {
		t.Error("expected regex match")
	}
	if r.Match("WARNING: something") {
		t.Error("expected no match for WARNING (not exact word WARN)")
	}
}

func TestRule_InvalidRegex(t *testing.T) {
	r := Rule{Regex: `[invalid`}
	if err := r.Compile(); err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestRule_Level(t *testing.T) {
	r := Rule{Level: "debug"}
	if err := r.Compile(); err != nil {
		t.Fatal(err)
	}
	if !r.Match("[DEBUG] verbose output") {
		t.Error("expected level match (case-insensitive)")
	}
	if r.Match("[INFO] normal output") {
		t.Error("expected no match for non-debug line")
	}
}

func TestChain_EmptyMatchesAll(t *testing.T) {
	c, err := NewChain(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Match("any line at all") {
		t.Error("empty chain should match every line")
	}
}

func TestChain_MultipleRules(t *testing.T) {
	rules := []Rule{
		{Contains: "ERROR"},
		{Contains: "database"},
	}
	c, err := NewChain(rules)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Match("ERROR connecting to database") {
		t.Error("expected match when all rules satisfied")
	}
	if c.Match("ERROR reading file") {
		t.Error("expected no match when only one rule satisfied")
	}
}
