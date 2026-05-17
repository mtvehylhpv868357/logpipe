package labelinject

import (
	"strings"
	"testing"
)

func TestNew_EmptyMapReturnsInjector(t *testing.T) {
	inj := New(map[string]string{})
	if inj == nil {
		t.Fatal("expected non-nil injector")
	}
}

func TestTransform_NoLabels_PassesThrough(t *testing.T) {
	inj := New(map[string]string{})
	line := "hello world"
	if got := inj.Transform(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestTransform_SingleLabel(t *testing.T) {
	inj := New(map[string]string{"env": "prod"})
	got := inj.Transform("msg")
	if !strings.Contains(got, "env=prod") {
		t.Errorf("expected label in output, got %q", got)
	}
	if !strings.HasPrefix(got, "msg") {
		t.Errorf("original line should be preserved as prefix, got %q", got)
	}
}

func TestTransform_MultipleLabels(t *testing.T) {
	inj := New(map[string]string{"env": "staging", "region": "us-east"})
	got := inj.Transform("log line")
	if !strings.Contains(got, "env=staging") {
		t.Errorf("missing env label in %q", got)
	}
	if !strings.Contains(got, "region=us-east") {
		t.Errorf("missing region label in %q", got)
	}
}

func TestNew_SkipsEmptyKey(t *testing.T) {
	inj := New(map[string]string{"": "value", "k": "v"})
	got := inj.Transform("line")
	if strings.Contains(got, "=value") {
		t.Errorf("should not inject label with empty key, got %q", got)
	}
	if !strings.Contains(got, "k=v") {
		t.Errorf("expected k=v in output, got %q", got)
	}
}

func TestNew_SkipsEmptyValue(t *testing.T) {
	inj := New(map[string]string{"k": "", "env": "dev"})
	got := inj.Transform("line")
	if strings.Contains(got, "k=") {
		t.Errorf("should not inject label with empty value, got %q", got)
	}
	if !strings.Contains(got, "env=dev") {
		t.Errorf("expected env=dev in output, got %q", got)
	}
}

func TestTransform_PreservesOriginalLine(t *testing.T) {
	inj := New(map[string]string{"host": "box1"})
	orig := "level=info msg=started"
	got := inj.Transform(orig)
	if !strings.HasPrefix(got, orig) {
		t.Errorf("original line must be a prefix of output, got %q", got)
	}
}
