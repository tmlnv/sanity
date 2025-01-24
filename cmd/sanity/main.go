package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/tui"
)

var version = "v1.0.0"

func main() {
	cfg := parseFlags()

	if !cfg.FlagsProvided {
		p := tea.NewProgram(tui.InitialModel())
		if _, err := p.Run(); err != nil {
			logger.Error("Error running TUI", "error", err)
			os.Exit(1)
		}
		return
	}

	runCLI(cfg)
}

func parseFlags() config.Config {
	var cfg config.Config

	flag.StringVar(&cfg.Prefix, "prefix", "", "Vanity prefix for Solana address")
	flag.StringVar(&cfg.Suffix, "suffix", "", "Vanity suffix for Solana address")
	flag.StringVar(&cfg.Regex, "regex", "", "Regex pattern to match")
	flag.IntVar(&cfg.NumAddresses, "count", 1, "Number of addresses to find (0=infinite)")
	flag.IntVar(&cfg.Concurrency, "threads", 0, "Number of worker threads (0=auto)")
	flag.DurationVar(&cfg.Timeout, "timeout", 0, "Maximum search duration")
	flag.StringVar(&cfg.LogFile, "logfile", "", "Path to log file")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("sanity %s\n", version)
		os.Exit(0)
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.NumCPU()
	}

	cfg.FlagsProvided = flag.NFlag() > 0
	return cfg
}

func runCLI(cfg config.Config) {
	logger.Init(cfg.LogFile)
	ctx, cancel := createContext(cfg)
	defer cancel()

	logger.Info("Starting vanity generation",
		"prefix", cfg.Prefix,
		"suffix", cfg.Suffix,
		"regex", cfg.Regex,
		"threads", cfg.Concurrency,
		"timeout", cfg.Timeout,
	)

	generator.Start(ctx, cancel, cfg, nil, false)
}

func createContext(cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Timeout > 0 {
		return context.WithTimeout(context.Background(), cfg.Timeout)
	}
	return context.WithCancel(context.Background())
}
