package router

import (
	"testing"

	"github.com/user/logpipe/config"
)

func multilineOutput(start, cont string, timeoutSec int, join string) config.Output {
	return config.Output{
		Type: "stdout",
		Multiline: &config.MultilineConfig{
			StartPattern:    start,
			ContinuePattern: cont,
			TimeoutSeconds:  timeoutSec,
			Join:            join,
		},
	}
}

func TestMultilineAgg_NilWhenNoConfig(t *testing.T) {
	o := config.Output{Type: "stdout"}
	if multilineAggForOutput(o) != nil {
		t.Fatal("expected nil when multiline config absent")
	}
}

func TestMultilineAgg_NilWhenEmptyStartPattern(t *testing.T) {
	o := multilineOutput("", "", 0, "")
	if multilineAggForOutput(o) != nil {
		t.Fatal("expected nil for empty start pattern")
	}
}

func TestMultilineAgg_ReturnsAggregator(t *testing.T) {
	o := multilineOutput("^ERROR", `^\s`, 5, "\n")
	agg := multilineAggForOutput(o)
	if agg == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestMultilineAgg_NilOnInvalidPattern(t *testing.T) {
	o := multilineOutput("[", "", 0, "")
	if multilineAggForOutput(o) != nil {
		t.Fatal("expected nil for invalid regex")
	}
}

func TestMultilineAgg_NilContinuePatternAllowed(t *testing.T) {
	o := multilineOutput("^INFO", "", 0, "")
	agg := multilineAggForOutput(o)
	if agg == nil {
		t.Fatal("expected aggregator with empty continue pattern")
	}
}
