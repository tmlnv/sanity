package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/tui"
)

var Version = "1.0.0"

func main() {
	cfg := parseFlags()

	if !cfg.FlagsProvided {
		startInteractiveTUI()
		return
	}

	initializeSystem(cfg)
}

func parseFlags() config.Config {
	var cfg config.Config

	flag.StringVar(&cfg.Prefix, "prefix", "", "Vanity prefix for Solana address")
	flag.StringVar(&cfg.Suffix, "suffix", "", "Vanity suffix for Solana address")
	flag.StringVar(&cfg.Regex, "regex", "", "Regex pattern to match")
	flag.IntVar(&cfg.NumAddresses, "count", 1, "Number of addresses to find (0=infinite)")
	flag.IntVar(&cfg.Concurrency, "threads", 0, "Number of worker threads (0=auto)")
	flag.DurationVar(&cfg.Timeout, "timeout", 0, "Maximum search duration")
	flag.StringVar(&cfg.LogFile, "logfile", "sanity.log", "Path to log file")
	versionFlag := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("sanity %s\n", Version)
		os.Exit(0)
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.NumCPU()
	}

	cfg.FlagsProvided = flag.NFlag() > 0
	return cfg
}

func startInteractiveTUI() {
	program := tui.NewProgram(tui.InitialModel())
	if _, err := program.Run(); err != nil {
		logger.Error("TUI failed", "error", err)
		os.Exit(1)
	}
}

func initializeSystem(cfg config.Config) {
	logger.Init(cfg.LogFile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Timeout > 0 {
		var timeoutCancel context.CancelFunc
		ctx, timeoutCancel = context.WithTimeout(ctx, cfg.Timeout)
		defer timeoutCancel()
	}

	logger.Info("Starting vanity generation",
		"prefix", cfg.Prefix,
		"suffix", cfg.Suffix,
		"regex", cfg.Regex,
		"threads", cfg.Concurrency,
		"timeout", cfg.Timeout,
	)

	generator.Start(ctx, cfg)
}
