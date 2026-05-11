package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/logpipe/internal/config"
	"github.com/user/logpipe/internal/pipeline"
)

const version = "0.1.0"

func main() {
	var (
		configFile  = flag.String("config", "logpipe.yaml", "path to configuration file")
		showVersion = flag.Bool("version", false, "print version and exit")
		verbose     = flag.Bool("verbose", false, "enable verbose logging")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("logpipe version %s\n", version)
		os.Exit(0)
	}

	// Load configuration from file
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		cfg.Verbose = true
	}

	// Build and run the pipeline
	p, err := pipeline.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building pipeline: %v\n", err)
		os.Exit(1)
	}

	// Read from stdin and process until EOF or error
	if err := p.Run(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "pipeline error: %v\n", err)
		os.Exit(1)
	}
}
