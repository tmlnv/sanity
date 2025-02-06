package tui

import (
	"context"
	"fmt"
	"runtime"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/ctx"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/validator"
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
		inputs: make([]textinput.Model, 6),
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
		{1, "Vanity suffix:", nil},
		{2, "Vanity regexp:", nil},
		{3, "Number of addresses to find (0=infinite):", validator.ValidateNumber},
		{4, fmt.Sprintf("Threads (0=auto, CPUs available: %d):", runtime.NumCPU()), validator.ValidateNumber},
		{5, "Timeout (e.g. 30s, 5m):", validator.ValidateDuration},
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
