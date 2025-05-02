// internal/config/config.go
package config

import (
	"strconv"
	"time"
)

type Config struct {
	Prefix          string
	Suffix          string
	Regexp          string
	NumAddresses    int
	Concurrency     int
	Timeout         time.Duration
	LogInterval     time.Duration
	LogFile         string
	PrivateKeysFile string
	FlagsProvided   bool
}

func (c Config) String() string {
	// Not including LogInterval in the string representation, as it's not relevant for TUI
	return "Prefix: " + c.Prefix + "\n" +
		"Suffix: " + c.Suffix + "\n" +
		"Regexp: " + c.Regexp + "\n" +
		"NumAddresses: " + strconv.Itoa(c.NumAddresses) + "\n" +
		"Concurrency: " + strconv.Itoa(c.Concurrency) + "\n" +
		"Timeout: " + c.Timeout.String() + "\n" +
		"LogInterval: " + c.LogInterval.String() + "\n" +
		"LogFile: " + c.LogFile + "\n" +
		"PrivateKeysFile: " + c.PrivateKeysFile + "\n"
}
