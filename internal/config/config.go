// internal/config/config.go
package config

import (
	"strconv"
	"time"
)

type Config struct {
	Prefix          string
	Suffix          string
	Regex           string
	NumAddresses    int
	Concurrency     int
	Timeout         time.Duration
	LogFile         string
	PrivateKeysFile string
	FlagsProvided   bool
}

func (c Config) String() string {
	return "Prefix: " + c.Prefix + "\n" +
		"Suffix: " + c.Suffix + "\n" +
		"Regex: " + c.Regex + "\n" +
		"NumAddresses: " + strconv.Itoa(c.NumAddresses) + "\n" +
		"Concurrency: " + strconv.Itoa(c.Concurrency) + "\n" +
		"Timeout: " + c.Timeout.String() + "\n" +
		"LogFile: " + c.LogFile + "\n" +
		"PrivateKeysFile: " + c.PrivateKeysFile + "\n"
}
