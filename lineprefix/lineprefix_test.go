package lineprefix_test

import (
	"fmt"
	"testing"

	"logpipe/lineprefix"
)

func TestPrepender_EmptyPrefix_PassesThrough(t *testing.T) {
	p := lineprefix.New("")
	got := p.Transform("hello world")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestPrepender_AddsPrefix(t *testing.T) {
	p := lineprefix.New("[INFO] ")
	got := p.Transform("service started")
	want := "[INFO] service started"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPrepender_EmptyLine(t *testing.T) {
	p := lineprefix.New(">> ")
	got := p.Transform("")
	want := ">> "
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPrepender_PreservesWhitespace(t *testing.T) {
	p := lineprefix.New("X")
	got := p.Transform("  leading spaces")
	want := "X  leading spaces"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDynamicPrepender_CallsFn(t *testing.T) {
	calls := 0
	dp := lineprefix.NewDynamic(func() string {
		calls++
		return fmt.Sprintf("[call%d] ", calls)
	})

	first := dp.Transform("line one")
	second := dp.Transform("line two")

	if first != "[call1] line one" {
		t.Fatalf("unexpected first: %q", first)
	}
	if second != "[call2] line two" {
		t.Fatalf("unexpected second: %q", second)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDynamicPrepender_EmptyFnResult_PassesThrough(t *testing.T) {
	dp := lineprefix.NewDynamic(func() string { return "" })
	got := dp.Transform("unchanged")
	if got != "unchanged" {
		t.Fatalf("expected %q, got %q", "unchanged", got)
	}
}

func TestNewDynamic_NilFn_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fn")
		}
	}()
	lineprefix.NewDynamic(nil)
}
