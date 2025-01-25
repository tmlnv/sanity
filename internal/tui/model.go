package tui

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/context"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
)

type Model struct {
	inputs       []textinput.Model
	step         int
	spinner      spinner.Model
	isGenerating bool
	isFinished   bool
	config       config.Config
	stats        generator.Stats
	lastResults  []string
	err          error
	updateChan   chan generator.StatsUpdate
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
	logger.Init(m.config.LogFile)
	ctx, cancel := context.CreateContext(m.config)
	defer cancel()
	defer func() { m.isFinished = true }()

	generator.Start(ctx, cancel, m.config, m.updateChan, true)
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
		return nil
	}
	_, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration format (e.g., 30s, 5m)")
	}
	return nil
}
