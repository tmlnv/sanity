package main

import (
	"os"

	"github.com/tmlnv/sanity/internal/cli"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/saver"
	"github.com/tmlnv/sanity/internal/tui"
)

func main() {
	cfg := cli.ParseFlags()
	logger.Init(cfg.LogFile)
	defer logger.Close()

	saver.Init(cfg.PrivateKeysFile)
	defer saver.Close()

	if !cfg.FlagsProvided {
		p := tui.NewProgram(tui.InitialModel())
		if _, err := p.Run(); err != nil {
			logger.Error("Error running TUI", "error", err)
			os.Exit(1)
		}
		return
	}

	cli.RunCLI(cfg)
}
