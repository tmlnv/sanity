package tui

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
)

type Model struct {
	step        int
	config      config.Config
	inputs      []textinput.Model
	spinner     spinner.Model
	generating  bool
	stats       generator.Stats
	lastResults []string
	err         error
	updateChan  chan generator.StatsUpdate
}

func InitialModel() Model {
	m := Model{
		config: config.Config{
			NumAddresses: 1,
			Concurrency:  runtime.NumCPU(),
		},
		inputs: make([]textinput.Model, 4),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("206"))),
		),
		lastResults: make([]string, 0),
		updateChan:  make(chan generator.StatsUpdate, 100),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32
		t.PromptStyle = t.PromptStyle.Faint(true)

		switch i {
		case 0:
			t.Prompt = "Enter vanity prefix: "
			t.Placeholder = "sol"
		case 1:
			t.Prompt = "Number of addresses to find (0=infinite): "
			t.Validate = ValidateNumber
		case 2:
			t.Prompt = fmt.Sprintf("Threads to use (%d available): ", runtime.NumCPU())
			t.Validate = ValidateNumber
		case 3:
			t.Prompt = "Timeout (e.g. 30s, 5m): "
			t.Validate = ValidateDuration
		}
		m.inputs[i] = t
	}
	return m
}

func (m *Model) startGeneration() {
	// Initialize logger before starting generator
	logger.Init(m.config.LogFile)

	ctx, cancel := context.WithCancel(context.Background())
	if m.config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, m.config.Timeout)
	}
	defer cancel()

	// Start generator with the update channel
	generator.Start(ctx, m.config, m.updateChan)
}

func ValidateNumber(s string) error {
	if s == "" || s == "0" {
		return nil
	}
	_, err := strconv.Atoi(s)
	return err
}

func ValidateDuration(s string) error {
	if s == "" {
		return nil
	}
	_, err := time.ParseDuration(s)
	return err
}
