package router_test

import (
	"bytes"
	"errors"
	"io"
	"logpipe/filter"
	"logpipe/router"
	"testing"
)

// bufWriteCloser wraps a bytes.Buffer as an io.WriteCloser.
type bufWriteCloser struct {
	buf    bytes.Buffer
	closed bool
}

func (b *bufWriteCloser) Write(p []byte) (int, error) { return b.buf.Write(p) }
func (b *bufWriteCloser) Close() error               { b.closed = true; return nil }

// errWriter always returns an error on Write.
type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) { return 0, errors.New("write error") }
func (e *errWriter) Close() error                { return nil }

func newPassChain() *filter.Chain { return filter.NewChain(nil) }

func TestDispatch_MatchingRoute(t *testing.T) {
	buf := &bufWriteCloser{}
	r := router.New([]*router.Route{
		{Chain: newPassChain(), Writer: buf},
	})
	n, err := r.Dispatch([]byte("hello world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 match, got %d", n)
	}
	if buf.buf.String() != "hello world\n" {
		t.Fatalf("unexpected output: %q", buf.buf.String())
	}
}

func TestDispatch_NoMatch(t *testing.T) {
	buf := &bufWriteCloser{}
	rules := []filter.Rule{{Field: "contains", Value: "ERROR"}}
	chain := filter.NewChain(rules)
	r := router.New([]*router.Route{
		{Chain: chain, Writer: buf},
	})
	n, err := r.Dispatch([]byte("info: all good"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Fatalf("expected 0 matches, got %d", n)
	}
}

func TestDispatch_WriteError(t *testing.T) {
	r := router.New([]*router.Route{
		{Chain: newPassChain(), Writer: &errWriter{}},
	})
	_, err := r.Dispatch([]byte("line"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClose_ClosesAllWriters(t *testing.T) {
	a, b := &bufWriteCloser{}, &bufWriteCloser{}
	r := router.New([]*router.Route{
		{Chain: newPassChain(), Writer: a},
		{Chain: newPassChain(), Writer: b},
	})
	_ = r.Close()
	if !a.closed || !b.closed {
		t.Fatal("expected both writers to be closed")
	}
}

// Ensure router.Route.Writer satisfies io.WriteCloser at compile time.
var _ io.WriteCloser = (*bufWriteCloser)(nil)
