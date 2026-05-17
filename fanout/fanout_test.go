package fanout_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"logpipe/fanout"
)

// bufCloser wraps a bytes.Buffer and implements io.WriteCloser.
type bufCloser struct {
	bytes.Buffer
	closed bool
	writeErr error
	closeErr error
}

func (b *bufCloser) Write(p []byte) (int, error) {
	if b.writeErr != nil {
		return 0, b.writeErr
	}
	return b.Buffer.Write(p)
}

func (b *bufCloser) Close() error {
	b.closed = true
	return b.closeErr
}

func TestFanout_WritesToAllWriters(t *testing.T) {
	a, b := &bufCloser{}, &bufCloser{}
	fw := fanout.New(a, b)

	_, err := fw.WriteString("hello\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.String() != "hello\n" {
		t.Errorf("writer a got %q, want %q", a.String(), "hello\n")
	}
	if b.String() != "hello\n" {
		t.Errorf("writer b got %q, want %q", b.String(), "hello\n")
	}
}

func TestFanout_ContinuesOnWriteError(t *testing.T) {
	failing := &bufCloser{writeErr: errors.New("disk full")}
	good := &bufCloser{}
	fw := fanout.New(failing, good)

	_, err := fw.WriteString("line\n")
	if err == nil {
		t.Fatal("expected error from failing writer")
	}
	// good writer should still receive the line
	if good.String() != "line\n" {
		t.Errorf("good writer got %q, want %q", good.String(), "line\n")
	}
}

func TestFanout_CloseClosesAll(t *testing.T) {
	a, b := &bufCloser{}, &bufCloser{}
	fw := fanout.New(a, b)

	if err := fw.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !a.closed {
		t.Error("writer a was not closed")
	}
	if !b.closed {
		t.Error("writer b was not closed")
	}
}

func TestFanout_CloseReturnsFirstError(t *testing.T) {
	a := &bufCloser{closeErr: errors.New("close failed")}
	b := &bufCloser{}
	fw := fanout.New(a, b)

	err := fw.Close()
	if err == nil {
		t.Fatal("expected close error")
	}
	if !b.closed {
		t.Error("writer b should still be closed despite error in a")
	}
}

func TestFanout_EmptyWriters(t *testing.T) {
	fw := fanout.New()
	if fw.Len() != 0 {
		t.Errorf("expected 0 writers, got %d", fw.Len())
	}
	_, err := fw.WriteString("anything")
	if err != nil {
		t.Errorf("unexpected error on empty fanout: %v", err)
	}
}

func TestFanout_ImplementsWriteCloser(t *testing.T) {
	var _ io.WriteCloser = fanout.New()
}
