package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/constants"
	"github.com/tmlnv/sanity/internal/ctx"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
)

func ParseFlags() config.Config {
	var cfg config.Config

	flag.StringVar(&cfg.Prefix, "prefix", "", "Vanity prefix for Solana address")
	flag.StringVar(&cfg.Suffix, "suffix", "", "Vanity suffix for Solana address")
	flag.StringVar(&cfg.Regex, "regex", "", "Regex pattern to match")
	flag.IntVar(&cfg.NumAddresses, "count", 1, "Number of addresses to find (0=infinite)")
	flag.IntVar(&cfg.Concurrency, "threads", 0, "Number of worker threads (0=auto)")
	flag.DurationVar(&cfg.Timeout, "timeout", 0, "Maximum search duration")
	flag.StringVar(&cfg.LogFile, "logfile", constants.LogFile, "Path to log file")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("sanity %s\n", constants.Version)
		os.Exit(0)
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.NumCPU()
	}

	cfg.FlagsProvided = flag.NFlag() > 0
	return cfg
}

func RunCLI(cfg config.Config) {
	logger.Init(cfg.LogFile)
	ctx, cancel := ctx.CreateContext(cfg)
	defer cancel()

	logger.Info("Starting vanity generation",
		"prefix", cfg.Prefix,
		"suffix", cfg.Suffix,
		"regex", cfg.Regex,
		"threads", cfg.Concurrency,
		"timeout", cfg.Timeout,
	)

	generator.Start(ctx, cfg, nil, false)
}
