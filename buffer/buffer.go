package buffer

import (
	"sync"
	"time"
)

// FlushFunc is called with a batch of lines when the buffer flushes.
type FlushFunc func(lines []string)

// Buffer accumulates lines and flushes them either when the batch
// reaches maxSize or when the flush interval elapses.
type Buffer struct {
	mu       sync.Mutex
	lines    []string
	maxSize  int
	interval time.Duration
	flushFn  FlushFunc
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// New creates a Buffer. maxSize <= 0 disables size-based flushing;
// interval <= 0 disables time-based flushing.
func New(maxSize int, interval time.Duration, fn FlushFunc) *Buffer {
	b := &Buffer{
		lines:   make([]string, 0, maxSize),
		maxSize: maxSize,
		interval: interval,
		flushFn: fn,
		stopCh:  make(chan struct{}),
	}
	if interval > 0 {
		b.wg.Add(1)
		go b.tickLoop()
	}
	return b
}

// Add appends a line to the buffer, flushing if maxSize is reached.
func (b *Buffer) Add(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, line)
	if b.maxSize > 0 && len(b.lines) >= b.maxSize {
		b.flush()
	}
}

// Flush forces an immediate flush of buffered lines.
func (b *Buffer) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flush()
}

// Stop halts the background ticker and performs a final flush.
func (b *Buffer) Stop() {
	close(b.stopCh)
	b.wg.Wait()
	b.Flush()
}

// flush drains lines and calls flushFn. Caller must hold b.mu.
func (b *Buffer) flush() {
	if len(b.lines) == 0 {
		return
	}
	batch := make([]string, len(b.lines))
	copy(batch, b.lines)
	b.lines = b.lines[:0]
	b.flushFn(batch)
}

func (b *Buffer) tickLoop() {
	defer b.wg.Done()
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			b.Flush()
		case <-b.stopCh:
			return
		}
	}
}
