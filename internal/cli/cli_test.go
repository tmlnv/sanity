package cli

import (
	"flag"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/logger"
)

func TestParseFlags(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	logger.Init("test.log")
	defer logger.Close()

	tests := []struct {
		name string
		args []string
		want config.Config
	}{
		{
			name: "default config",
			args: []string{"sanity"},
			want: config.Config{
				NumAddresses:    1,
				Concurrency:     runtime.NumCPU(),
				LogFile:         config.LogFile,
				PrivateKeysFile: config.PrivateKeysFile,
			},
		},
		{
			name: "config with prefix and suffix",
			args: []string{"sanity", "-prefix", "ABC", "-suffix", "XYZ"},
			want: config.Config{
				Prefix:          "ABC",
				Suffix:          "XYZ",
				NumAddresses:    1,
				Concurrency:     runtime.NumCPU(),
				LogFile:         config.LogFile,
				PrivateKeysFile: config.PrivateKeysFile,
				FlagsProvided:   true,
			},
		},
		{
			name: "config with regexp pattern",
			args: []string{"sanity", "-regexp", "^[A-Za-z0-9]{44}$"},
			want: config.Config{
				Regexp:          "^[A-Za-z0-9]{44}$",
				NumAddresses:    1,
				Concurrency:     runtime.NumCPU(),
				LogFile:         config.LogFile,
				PrivateKeysFile: config.PrivateKeysFile,
				FlagsProvided:   true,
			},
		},
		{
			name: "config with all fields",
			args: []string{"sanity", "-prefix", "ABC", "-suffix", "XYZ", "-regexp", "^ABC.*XYZ$"},
			want: config.Config{
				Prefix:          "ABC",
				Suffix:          "XYZ",
				Regexp:          "^ABC.*XYZ$",
				NumAddresses:    1,
				Concurrency:     runtime.NumCPU(),
				LogFile:         config.LogFile,
				PrivateKeysFile: config.PrivateKeysFile,
				FlagsProvided:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Set test args
			os.Args = tt.args

			got := ParseFlags()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunCLI(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "run with empty config",
			args: args{cfg: config.Config{}},
		},
		{
			name: "run with prefix and suffix",
			args: args{cfg: config.Config{
				Prefix: "ABC",
				Suffix: "XYZ",
			}},
		},
		{
			name: "run with regexp pattern",
			args: args{cfg: config.Config{
				Regexp: "^[A-Za-z0-9]{44}$",
			}},
		},
		{
			name: "run with all fields",
			args: args{cfg: config.Config{
				Prefix: "ABC",
				Suffix: "XYZ",
				Regexp: "^ABC.*XYZ$",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunCLI(tt.args.cfg)
		})
	}
}
