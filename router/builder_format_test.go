package router

import (
	"testing"

	"github.com/yourorg/logpipe/config"
)

func formatOutput(fmtType, template string) config.Output {
	return config.Output{
		Type: "stdout",
		Format: &config.FormatConfig{
			Type:     fmtType,
			Template: template,
		},
	}
}

func TestFormatterForOutput_NilWhenNoFormat(t *testing.T) {
	out := config.Output{Type: "stdout"}
	if f := formatterForOutput(out); f != nil {
		t.Errorf("expected nil formatter, got %T", f)
	}
}

func TestFormatterForOutput_JSON(t *testing.T) {
	out := formatOutput("json", "")
	f := formatterForOutput(out)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	result := f.Format(map[string]string{"k": "v"})
	if result != `{"k":"v"}` {
		t.Errorf("unexpected json output: %s", result)
	}
}

func TestFormatterForOutput_Logfmt(t *testing.T) {
	out := formatOutput("logfmt", "")
	f := formatterForOutput(out)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	result := f.Format(map[string]string{"level": "info"})
	if result != "level=info" {
		t.Errorf("unexpected logfmt output: %s", result)
	}
}

func TestFormatterForOutput_Template(t *testing.T) {
	out := formatOutput("template", "{{level}}: {{msg}}")
	f := formatterForOutput(out)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	result := f.Format(map[string]string{"level": "debug", "msg": "ok"})
	if result != "debug: ok" {
		t.Errorf("unexpected template output: %s", result)
	}
}

func TestFormatterForOutput_TemplateEmptyTemplate(t *testing.T) {
	out := formatOutput("template", "")
	if f := formatterForOutput(out); f != nil {
		t.Errorf("expected nil for empty template, got %T", f)
	}
}

func TestFormatterForOutput_UnknownType(t *testing.T) {
	out := formatOutput("xml", "")
	if f := formatterForOutput(out); f != nil {
		t.Errorf("expected nil for unknown type, got %T", f)
	}
}
