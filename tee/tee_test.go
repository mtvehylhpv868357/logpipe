package tee_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"logpipe/tee"
)

// closableBuffer wraps bytes.Buffer and satisfies io.WriteCloser.
type closableBuffer struct {
	bytes.Buffer
	closed bool
}

func (c *closableBuffer) Close() error {
	c.closed = true
	return nil
}

// errorWriter always returns an error on Write.
type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("secondary error")
}

func TestTee_WritesToPrimary(t *testing.T) {
	primary := &closableBuffer{}
	w := tee.New(primary, nil)

	_, err := w.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := primary.String(); got != "hello\n" {
		t.Errorf("primary = %q, want %q", got, "hello\n")
	}
}

func TestTee_CopiestoSecondary(t *testing.T) {
	primary := &closableBuffer{}
	var secondary bytes.Buffer
	w := tee.New(primary, &secondary)

	_, err := w.Write([]byte("world\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := secondary.String(); got != "world\n" {
		t.Errorf("secondary = %q, want %q", got, "world\n")
	}
}

func TestTee_SecondaryErrorDoesNotAffectPrimary(t *testing.T) {
	primary := &closableBuffer{}
	w := tee.New(primary, &errorWriter{})

	_, err := w.Write([]byte("line\n"))
	if err != nil {
		t.Fatalf("unexpected error from primary: %v", err)
	}
	if got := primary.String(); got != "line\n" {
		t.Errorf("primary = %q, want %q", got, "line\n")
	}
}

func TestTee_NilSecondaryIsNoop(t *testing.T) {
	primary := &closableBuffer{}
	w := tee.New(primary, nil)

	_, err := w.Write([]byte("noop\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTee_CloseClosesPrimary(t *testing.T) {
	primary := &closableBuffer{}
	w := tee.New(primary, nil)

	if err := w.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !primary.closed {
		t.Error("expected primary to be closed")
	}
}

func TestTee_ImplementsWriteCloser(t *testing.T) {
	primary := &closableBuffer{}
	var _ io.WriteCloser = tee.New(primary, nil)
}
