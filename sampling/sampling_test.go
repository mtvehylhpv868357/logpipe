package sampling_test

import (
	"math/rand"
	"testing"

	"logpipe/sampling"
)

func TestRateSampler_AlwaysForwardsAtRate1(t *testing.T) {
	s := sampling.NewRateSampler(1.0, rand.NewSource(42))
	for i := 0; i < 100; i++ {
		if !s.Sample("line") {
			t.Fatal("expected sample=true at rate 1.0")
		}
	}
}

func TestRateSampler_NeverForwardsAtRate0(t *testing.T) {
	s := sampling.NewRateSampler(0.0, rand.NewSource(42))
	for i := 0; i < 100; i++ {
		if s.Sample("line") {
			t.Fatal("expected sample=false at rate 0.0")
		}
	}
}

func TestRateSampler_ClampsNegativeRate(t *testing.T) {
	s := sampling.NewRateSampler(-5.0, rand.NewSource(1))
	for i := 0; i < 50; i++ {
		if s.Sample("x") {
			t.Fatal("negative rate should be clamped to 0")
		}
	}
}

func TestRateSampler_ClampsRateAbove1(t *testing.T) {
	s := sampling.NewRateSampler(99.0, rand.NewSource(1))
	for i := 0; i < 50; i++ {
		if !s.Sample("x") {
			t.Fatal("rate above 1 should be clamped to 1.0")
		}
	}
}

func TestRateSampler_ApproximateRate(t *testing.T) {
	s := sampling.NewRateSampler(0.5, rand.NewSource(99))
	hits := 0
	const n = 10000
	for i := 0; i < n; i++ {
		if s.Sample("line") {
			hits++
		}
	}
	ratio := float64(hits) / float64(n)
	if ratio < 0.45 || ratio > 0.55 {
		t.Errorf("expected ~0.5 hit rate, got %.3f", ratio)
	}
}

func TestNthSampler_ForwardsEveryNth(t *testing.T) {
	s := sampling.NewNthSampler(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Sample("line")
	}
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestNthSampler_ZeroBecomesOne(t *testing.T) {
	s := sampling.NewNthSampler(0)
	for i := 0; i < 10; i++ {
		if !s.Sample("x") {
			t.Fatal("n=0 should behave as n=1, forwarding all lines")
		}
	}
}

func TestPassthroughSampler_AlwaysTrue(t *testing.T) {
	var s sampling.PassthroughSampler
	for i := 0; i < 20; i++ {
		if !s.Sample("anything") {
			t.Fatal("PassthroughSampler must always return true")
		}
	}
}
