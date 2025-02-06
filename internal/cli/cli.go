package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/ctx"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/validator"
)

func ParseFlags() config.Config {
	var cfg config.Config
	var timeoutStr string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Sanity - Solana Vanity Address Generator\n\n")
		fmt.Fprintf(os.Stderr, "Usage: sanity [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  Generate an address with prefix 'sol':\n    $ sanity -prefix sol\n\n")
		fmt.Fprintf(os.Stderr, "  Generate 5 addresses with suffix 'eth':\n    $ sanity -suffix eth -count 5\n\n")
		fmt.Fprintf(os.Stderr, "  Generate addresses matching regex with timeout:\n    $ sanity -regex '^sol.*eth$' -timeout 5m\n")
	}

	flag.StringVar(&cfg.Prefix, "prefix", "", "Vanity prefix for Solana address")
	flag.StringVar(&cfg.Suffix, "suffix", "", "Vanity suffix for Solana address")
	flag.StringVar(&cfg.Regex, "regex", "", "Regex pattern to match")
	flag.IntVar(&cfg.NumAddresses, "count", 1, "Number of addresses to find (0=infinite)")
	flag.IntVar(&cfg.Concurrency, "threads", 0, "Number of worker threads (0=auto)")
	flag.StringVar(&timeoutStr, "timeout", "0", "Maximum search duration (e.g., 30s, 5m, or number of seconds)")
	flag.StringVar(&cfg.LogFile, "logfile", config.LogFile, "Path to log file")
	flag.StringVar(&cfg.PrivateKeysFile, "private-keys", config.PrivateKeysFile, "Path to private keys file")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	if *showVersion {
		fmt.Printf("sanity %s\n", config.Version)
		os.Exit(0)
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.NumCPU()
	}

	// Validate and parse timeout
	if duration, err := validator.ValidateTimeout(timeoutStr); err != nil {
		// fmt.Printf instead of logger.Error because logger is not initialized yet
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else {
		cfg.Timeout = duration
	}

	cfg.FlagsProvided = flag.NFlag() > 0
	return cfg
}

func RunCLI(cfg config.Config) {
	// Validate Solana address pattern before starting
	if err := validator.ValidateSolana(cfg.Prefix, cfg.Suffix, cfg.Regex); err != nil {
		logger.Error("Validation error", "error", err)
		os.Exit(1)
	}

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
