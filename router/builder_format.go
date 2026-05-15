package router

import (
	"github.com/yourorg/logpipe/config"
	"github.com/yourorg/logpipe/format"
)

// formatterForOutput returns a Formatter based on the output config, or nil
// if no format section is defined. Supported types: json, logfmt, template.
func formatterForOutput(out config.Output) format.Formatter {
	if out.Format == nil {
		return nil
	}
	switch out.Format.Type {
	case "json":
		return format.NewJSONFormatter()
	case "logfmt":
		return format.NewLogfmtFormatter()
	case "template":
		tmpl := out.Format.Template
		if tmpl == "" {
			return nil
		}
		return format.NewTemplateFormatter(tmpl)
	default:
		return nil
	}
}
