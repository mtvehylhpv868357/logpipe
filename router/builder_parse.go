package router

import (
	"github.com/user/logpipe/config"
	"github.com/user/logpipe/parse"
)

// parserForOutput returns a Parser for the output's configured format,
// or nil if no format is specified (passthrough behaviour).
func parserForOutput(out config.Output) parse.Parser {
	if out.Parse == "" {
		return nil
	}
	return parse.NewParser(out.Parse)
}
