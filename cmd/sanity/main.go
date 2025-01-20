package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/tui"
)

func main() {
	// Parse CLI flags
	prefix := flag.String("prefix", "", "Desired prefix for the Solana address")
	suffix := flag.String("suffix", "", "Desired suffix for the Solana address")
	regex := flag.String("regex", "", "Desired regex pattern for the Solana address")
	numWallets := flag.Int("wallets", 1, "Number of wallets to find (0 for infinite)")
	timeout := flag.Duration("timeout", 0, "Timeout for the generation process (e.g., 10s, 5m)")
	threads := flag.Int("threads", 0, "Number of threads to use (0 for automatic)")
	logFile := flag.String("logfile", "results.log", "File to log results to")
	quiet := flag.Bool("quiet", false, "Disable TUI and only log to file")
	flag.Parse()

	// Validate input
	if *prefix == "" && *suffix == "" && *regex == "" {
		fmt.Println("Please specify a prefix, suffix, or regex pattern.")
		flag.Usage()
		os.Exit(1)
	}

	// Set up configuration
	cfg := config.Config{
		Prefix:     *prefix,
		Suffix:     *suffix,
		Regex:      *regex,
		NumWallets: *numWallets,
		Timeout:    *timeout,
		Threads:    *threads,
		LogFile:    *logFile,
		Quiet:      *quiet,
	}

	// Initialize logger
	logger.InitLogger(cfg.LogFile, !cfg.Quiet)

	// Start TUI if not in quiet mode
	if !cfg.Quiet {
		go tui.StartTUI(&cfg)
	}

	// Start vanity address generation
	generator.StartGeneration(&cfg)
}
