package pipeline

import (
	"bufio"
	"io"
	"sync"

	"github.com/yourorg/logpipe/metrics"
	"github.com/yourorg/logpipe/router"
)

// Pipeline reads lines from a Reader, routes each line through the router,
// and tracks metrics. It is safe to run concurrently.
type Pipeline struct {
	reader  io.Reader
	router  *router.Router
	counters *metrics.Counters
}

// New creates a Pipeline that reads from r and routes via rt.
func New(r io.Reader, rt *router.Router, c *metrics.Counters) *Pipeline {
	return &Pipeline{
		reader:   r,
		router:   rt,
		counters: c,
	}
}

// Run scans lines from the reader until EOF or error, routing each line.
// It returns the first non-EOF read error encountered.
func (p *Pipeline) Run() error {
	scanner := bufio.NewScanner(p.reader)
	var wg sync.WaitGroup

	for scanner.Scan() {
		line := scanner.Text()
		p.counters.IncReceived()

		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			routed := p.router.Route(l)
			if routed {
				p.counters.IncRouted()
			} else {
				p.counters.IncDropped()
			}
		}(line)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
