package config

import "time"

// Config holds the configuration for the vanity address generator.
type Config struct {
	Prefix     string        // Desired prefix for the address
	Suffix     string        // Desired suffix for the address
	Regex      string        // Desired regex pattern for the address
	NumWallets int           // Number of wallets to find (0 for infinite)
	Timeout    time.Duration // Timeout for the generation process
	Threads    int           // Number of threads to use (0 for automatic)
	LogFile    string        // File to log results to
	Quiet      bool          // Disable TUI and only log to file
}
