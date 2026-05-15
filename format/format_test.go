package format_test

import (
	"testing"

	"github.com/yourorg/logpipe/format"
)

func TestJSONFormatter_EmptyFields(t *testing.T) {
	f := format.NewJSONFormatter()
	out := f.Format(map[string]string{})
	if out != "{}" {
		t.Errorf("expected {}, got %s", out)
	}
}

func TestJSONFormatter_SingleField(t *testing.T) {
	f := format.NewJSONFormatter()
	out := f.Format(map[string]string{"level": "info"})
	if out != `{"level":"info"}` {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestLogfmtFormatter_EmptyFields(t *testing.T) {
	f := format.NewLogfmtFormatter()
	out := f.Format(map[string]string{})
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestLogfmtFormatter_SortsKeys(t *testing.T) {
	f := format.NewLogfmtFormatter()
	out := f.Format(map[string]string{"z": "last", "a": "first"})
	if out != "a=first z=last" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestLogfmtFormatter_QuotesValuesWithSpaces(t *testing.T) {
	f := format.NewLogfmtFormatter()
	out := f.Format(map[string]string{"msg": "hello world"})
	if out != `msg="hello world"` {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateFormatter_Substitutes(t *testing.T) {
	f := format.NewTemplateFormatter("[{{level}}] {{msg}}")
	out := f.Format(map[string]string{"level": "warn", "msg": "disk full"})
	if out != "[warn] disk full" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateFormatter_MissingKeyLeftAsIs(t *testing.T) {
	f := format.NewTemplateFormatter("[{{level}}] {{msg}}")
	out := f.Format(map[string]string{"level": "error"})
	if out != "[error] {{msg}}" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestTemplateFormatter_EmptyTemplate(t *testing.T) {
	f := format.NewTemplateFormatter("")
	out := f.Format(map[string]string{"level": "info"})
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}
