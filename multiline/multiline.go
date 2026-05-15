// Package multiline provides a line aggregator that joins multi-line log
// entries (e.g. stack traces) into a single logical line.
package multiline

import (
	"regexp"
	"strings"
	"time"
)

// Aggregator buffers lines that belong to a single logical event and flushes
// them when a new event boundary is detected or a timeout elapses.
type Aggregator struct {
	start   *regexp.Regexp
	continue_ *regexp.Regexp
	timeout time.Duration
	join    string

	buf     []string
	lastAt  time.Time
	clock   func() time.Time
}

// Option configures an Aggregator.
type Option func(*Aggregator)

// WithJoin sets the string used to join buffered lines (default: "\n").
func WithJoin(sep string) Option {
	return func(a *Aggregator) { a.join = sep }
}

// WithTimeout sets the maximum idle time before a pending buffer is flushed.
func WithTimeout(d time.Duration) Option {
	return func(a *Aggregator) { a.timeout = d }
}

// New creates an Aggregator. startPattern marks the first line of a new event.
// If continuePattern is non-empty, only lines matching it are appended to the
// current event; any non-matching line is treated as a new event boundary.
func New(startPattern, continuePattern string, opts ...Option) (*Aggregator, error) {
	start, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, err
	}
	var cont *regexp.Regexp
	if continuePattern != "" {
		if cont, err = regexp.Compile(continuePattern); err != nil {
			return nil, err
		}
	}
	a := &Aggregator{
		start:     start,
		continue_: cont,
		timeout:   5 * time.Second,
		join:      "\n",
		clock:     time.Now,
	}
	for _, o := range opts {
		o(a)
	}
	return a, nil
}

// Add feeds a line to the aggregator. It returns a flushed event string and
// true when a complete event has been assembled, or "", false otherwise.
func (a *Aggregator) Add(line string) (string, bool) {
	now := a.clock()

	// Timeout flush: if the buffer has been idle too long, emit it first.
	timedOut := len(a.buf) > 0 && a.timeout > 0 && now.Sub(a.lastAt) >= a.timeout

	startsNew := a.start.MatchString(line)
	continues := !startsNew && a.continue_ != nil && a.continue_.MatchString(line)

	if timedOut || (startsNew && len(a.buf) > 0) {
		event := strings.Join(a.buf, a.join)
		a.buf = nil
		if startsNew {
			a.buf = []string{line}
			a.lastAt = now
		} else if continues {
			a.buf = []string{line}
			a.lastAt = now
		}
		return event, true
	}

	if startsNew || continues || a.continue_ == nil {
		a.buf = append(a.buf, line)
		a.lastAt = now
		return "", false
	}

	// Line does not continue the current event and does not start a new one;
	// flush pending buffer and emit this line immediately as its own event.
	if len(a.buf) > 0 {
		event := strings.Join(a.buf, a.join)
		a.buf = nil
		// We'll lose this line; push it as a standalone next call — return
		// the pending buffer and queue line for next Add.
		a.buf = []string{line}
		a.lastAt = now
		return event, true
	}

	return line, true
}

// Flush returns any buffered lines as a final event, clearing the buffer.
func (a *Aggregator) Flush() (string, bool) {
	if len(a.buf) == 0 {
		return "", false
	}
	event := strings.Join(a.buf, a.join)
	a.buf = nil
	return event, true
}
