package router

import (
	"io"
	"logpipe/filter"
)

// Route defines a single routing rule: a filter chain paired with a writer.
type Route struct {
	Chain  *filter.Chain
	Writer io.WriteCloser
}

// Router holds multiple routes and dispatches log lines to matching outputs.
type Router struct {
	routes []*Route
}

// New creates a new Router with the given routes.
func New(routes []*Route) *Router {
	return &Router{routes: routes}
}

// Dispatch sends line to every route whose filter chain matches.
// Returns the number of routes the line was written to and the first error
// encountered (remaining routes are still attempted).
func (r *Router) Dispatch(line []byte) (int, error) {
	var firstErr error
	matched := 0

	for _, route := range r.routes {
		if route.Chain.Match(string(line)) {
			if _, err := route.Writer.Write(append(line, '\n')); err != nil && firstErr == nil {
				firstErr = err
			}
			matched++
		}
	}
	return matched, firstErr
}

// Close closes all writers attached to routes.
func (r *Router) Close() error {
	var firstErr error
	for _, route := range r.routes {
		if err := route.Writer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
