// internal/config/config.go
package config

import "time"

type Config struct {
	Prefix        string
	Suffix        string
	Regex         string
	NumAddresses  int
	Concurrency   int
	Timeout       time.Duration
	LogFile       string
	FlagsProvided bool
}
