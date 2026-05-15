package router

import (
	"time"

	"github.com/yourorg/logpipe/buffer"
	"github.com/yourorg/logpipe/config"
	"github.com/yourorg/logpipe/output"
)

// bufferWrapWriter wraps an output.Writer with a Buffer so writes are
// batched before being forwarded to the underlying writer.
type bufferWrapWriter struct {
	buf *buffer.Buffer
}

func (bw *bufferWrapWriter) Write(line string) error {
	bw.buf.Add(line)
	return nil
}

func (bw *bufferWrapWriter) Close() error {
	bw.buf.Stop()
	return nil
}

// wrapWithBuffer optionally wraps w in a Buffer based on cfg.
// Returns w unchanged when no buffer config is present.
func wrapWithBuffer(w output.Writer, cfg *config.OutputConfig) output.Writer {
	bc := cfg.Buffer
	if bc == nil {
		return w
	}
	if bc.MaxSize <= 0 && bc.Interval == "" {
		return w
	}

	var interval time.Duration
	if bc.Interval != "" {
		var err error
		interval, err = time.ParseDuration(bc.Interval)
		if err != nil {
			interval = time.Second
		}
	}

	underlying := w
	flushFn := func(lines []string) {
		for _, l := range lines {
			_ = underlying.Write(l)
		}
	}

	buf := buffer.New(bc.MaxSize, interval, flushFn)
	return &bufferWrapWriter{buf: buf}
}
