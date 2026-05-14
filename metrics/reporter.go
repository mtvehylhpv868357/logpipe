package metrics

import (
	"io"
	"time"
)

// Reporter periodically prints metrics snapshots to a writer.
type Reporter struct {
	counters  *Counters
	interval  time.Duration
	out       io.Writer
	stopCh    chan struct{}
	doneCh    chan struct{}
}

// NewReporter creates a Reporter that writes to out every interval.
func NewReporter(c *Counters, interval time.Duration, out io.Writer) *Reporter {
	return &Reporter{
		counters: c,
		interval: interval,
		out:      out,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
}

// Start begins the periodic reporting loop in a goroutine.
func (r *Reporter) Start() {
	go func() {
		defer close(r.doneCh)
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.counters.TakeSnapshot().Print(r.out)
			case <-r.stopCh:
				// Print final snapshot on shutdown.
				r.counters.TakeSnapshot().Print(r.out)
				return
			}
		}
	}()
}

// Stop signals the reporter to cease and waits for it to finish.
func (r *Reporter) Stop() {
	close(r.stopCh)
	<-r.doneCh
}
