package main

import (
	"fmt"
	"os"

	"github.com/terminalizcrazy/terminalizcrazy/internal/config"
	"github.com/terminalizcrazy/terminalizcrazy/internal/tui"
)

// Version info - set via ldflags during build
var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Start the TUI application
	if err := tui.Run(cfg, version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
