package router_test

import (
	"logpipe/config"
	"logpipe/router"
	"os"
	"path/filepath"
	"testing"
)

func TestBuild_StdoutOutput(t *testing.T) {
	cfg := &config.Config{
		Outputs: []config.Output{
			{Type: "stdout"},
		},
	}
	r, err := router.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil router")
	}
	_ = r.Close()
}

func TestBuild_FileOutput(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.log")
	cfg := &config.Config{
		Outputs: []config.Output{
			{Type: "file", Path: tmp},
		},
	}
	r, err := router.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _ = r.Dispatch([]byte("test line"))
	_ = r.Close()

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if string(data) != "test line\n" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
}

func TestBuild_UnknownType(t *testing.T) {
	cfg := &config.Config{
		Outputs: []config.Output{
			{Type: "kafka"},
		},
	}
	_, err := router.Build(cfg)
	if err == nil {
		t.Fatal("expected error for unknown output type")
	}
}

func TestBuild_FileMissingPath(t *testing.T) {
	cfg := &config.Config{
		Outputs: []config.Output{
			{Type: "file"},
		},
	}
	_, err := router.Build(cfg)
	if err == nil {
		t.Fatal("expected error when file path is empty")
	}
}

func TestBuild_MultipleOutputs(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "multi.log")
	cfg := &config.Config{
		Outputs: []config.Output{
			{Type: "stdout"},
			{Type: "file", Path: tmp},
		},
	}
	r, err := router.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _ = r.Dispatch([]byte("multi output"))
	_ = r.Close()

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if string(data) != "multi output\n" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
}
