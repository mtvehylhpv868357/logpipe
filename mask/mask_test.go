package mask

import (
	"testing"
)

func TestFieldMasker_NoFields_PassesThrough(t *testing.T) {
	fm := New(nil)
	line := "user=alice token=secret123"
	if got := fm.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestFieldMasker_MasksNamedField(t *testing.T) {
	fm := New([]string{"token"})
	got := fm.Apply("user=alice token=secret123")
	want := "user=alice token=***"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_MasksQuotedValue(t *testing.T) {
	fm := New([]string{"password"})
	got := fm.Apply(`user=bob password="my secret"`)
	want := `user=bob password="***"`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_CaseInsensitiveKey(t *testing.T) {
	fm := New([]string{"Token"})
	got := fm.Apply("TOKEN=abc user=dave")
	want := "TOKEN=*** user=dave"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_MultipleFields(t *testing.T) {
	fm := New([]string{"token", "password"})
	got := fm.Apply("user=alice token=xyz password=hunter2 level=info")
	want := "user=alice token=*** password=*** level=info"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_CustomMask(t *testing.T) {
	fm := NewWithMask([]string{"secret"}, "[REDACTED]")
	got := fm.Apply("secret=abc123")
	want := "secret=[REDACTED]"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_EmptyMaskFallsBackToDefault(t *testing.T) {
	fm := NewWithMask([]string{"key"}, "")
	got := fm.Apply("key=value")
	want := "key=***"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldMasker_NoMatchPassesThrough(t *testing.T) {
	fm := New([]string{"secret"})
	line := "user=alice level=info msg=hello"
	if got := fm.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}
