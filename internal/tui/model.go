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
	"github.com/tmlnv/sanity/internal/ctx"
	"github.com/tmlnv/sanity/internal/generator"
)

const (
	isConfig = iota
	isGenerating
	isFinished
)

type modelState int

type Model struct {
	inputs        []textinput.Model
	step          int
	spinner       spinner.Model
	state         modelState
	config        config.Config
	stats         generator.Stats
	lastGenerated []string
	matched       []string
	err           error
	updateChan    chan generator.StatsUpdate
	cancel        context.CancelFunc
}

func InitialModel() Model {
	m := Model{
		config: config.Config{
			NumAddresses:    1,
			Concurrency:     runtime.NumCPU(),
			LogFile:         config.LogFile,
			PrivateKeysFile: config.PrivateKeysFile,
		},
		inputs: make([]textinput.Model, 4),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("206"))),
		),
		updateChan: make(chan generator.StatsUpdate, 100),
	}

	inputs := []struct {
		index      int
		prompt     string
		validation func(string) error
	}{
		{0, "Vanity prefix (e.g. 'sol'):", nil},
		{1, "Number of addresses to find (0=infinite):", validateNumber},
		{2, fmt.Sprintf("Threads (0=auto, CPUs available: %d):", runtime.NumCPU()), validateNumber},
		{3, "Timeout (e.g. 30s, 5m):", validateDuration},
	}

	for _, in := range inputs {
		ti := textinput.New()
		ti.Prompt = in.prompt
		ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if in.validation != nil {
			ti.Validate = in.validation
		}
		m.inputs[in.index] = ti
	}

	m.inputs[0].Focus()
	return m
}

func (m *Model) startGeneration() {
	m.state = isGenerating

	ctx, cancel := ctx.CreateContext(m.config)

	// Store cancel function in model to call later
	m.cancel = cancel

	go generator.Start(ctx, m.config, m.updateChan, true)
}

func validateNumber(s string) error {
	if s == "" || s == "0" {
		return nil
	}
	_, err := strconv.Atoi(s)
	return err
}

func validateDuration(s string) error {
	if s == "" {
		return nil // Allow empty input during typing
	}

	// Special case: allow standalone zero
	if s == "0" {
		return nil
	}

	// Allow pure numeric input (temporary during typing)
	if _, err := strconv.Atoi(s); err == nil {
		return nil
	}

	// Validate proper duration format
	_, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration format (e.g., 30s, 5m)")
	}
	return nil
}
