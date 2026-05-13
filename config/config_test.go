package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return tmp
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
outputs:
  - name: console
    type: stdout
  - name: errors
    type: file
    target: /var/log/errors.log
    filters:
      - field: level
        match: error
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(cfg.Outputs))
	}
	if cfg.Outputs[1].Filters[0].Match != "error" {
		t.Errorf("expected filter match 'error', got %q", cfg.Outputs[1].Filters[0].Match)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoOutputs(t *testing.T) {
	path := writeTemp(t, "outputs: []\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty outputs")
	}
}

func TestLoad_InvalidType(t *testing.T) {
	yaml := `
outputs:
  - name: bad
    type: kafka
    target: localhost:9092
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for unsupported output type")
	}
}

func TestLoad_MissingTargetForFile(t *testing.T) {
	yaml := `
outputs:
  - name: myfile
    type: file
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing target on file output")
	}
}
