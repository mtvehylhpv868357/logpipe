package router

import (
	"github.com/yourorg/logpipe/config"
	"github.com/yourorg/logpipe/labelinject"
)

// labelInjectorForOutput returns a *labelinject.Injector if the output config
// contains any static labels, or nil if none are defined.
func labelInjectorForOutput(o config.Output) *labelinject.Injector {
	if len(o.Labels) == 0 {
		return nil
	}
	return labelinject.New(o.Labels)
}
