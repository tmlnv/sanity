package config

import (
	"strings"
	"testing"
	"time"
)

func TestConfig_String(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want []string // Substrings that should be present in the output
	}{
		{
			name: "basic config",
			cfg: Config{
				Prefix:          "test",
				Suffix:          "suffix",
				Regexp:          "^[a-z]+$",
				NumAddresses:    10,
				Concurrency:     4,
				Timeout:         5 * time.Second,
				LogFile:         "test.log",
				PrivateKeysFile: "private.log",
				FlagsProvided:   true,
			},
			want: []string{
				"Prefix: test",
				"Suffix: suffix",
				"Regexp: ^[a-z]+$",
				"NumAddresses: 10",
				"Concurrency: 4",
				"Timeout: 5s",
				"LogInterval: 0s",
				"LogFile: test.log",
				"PrivateKeysFile: private.log",
			},
		},
		{
			name: "empty config",
			cfg:  Config{},
			want: []string{
				"Prefix: ",
				"Suffix: ",
				"Regexp: ",
				"NumAddresses: 0",
				"Concurrency: 0",
				"Timeout: 0s",
				"LogInterval: 0s",
				"LogFile: ",
				"PrivateKeysFile: ",
			},
		},
		{
			name: "config with log interval",
			cfg: Config{
				LogInterval: 10 * time.Minute,
			},
			want: []string{
				"LogInterval: 10m0s",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.String()
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("Config.String() = %v, should contain %v", got, want)
				}
			}
		})
	}
}
