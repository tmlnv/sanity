package tui

import (
	"context"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.step < len(m.inputs) {
				return m.handleInput()
			}
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
	return m, cmd
}

func (m Model) handleInput() (Model, tea.Cmd) {
	switch m.step {
	case 0:
		m.config.Prefix = m.inputs[0].Value()
	case 1:
		if val := m.inputs[1].Value(); val != "" {
			m.config.NumAddresses, _ = strconv.Atoi(val)
		}
	case 2:
		if val := m.inputs[2].Value(); val != "" {
			m.config.Concurrency, _ = strconv.Atoi(val)
		}
	case 3:
		if val := m.inputs[3].Value(); val != "" {
			m.config.Timeout, _ = time.ParseDuration(val)
		}
	}

	if m.step < len(m.inputs)-1 {
		m.step++
		return m, textinput.Blink
	}

	return m, tea.Quit
}

type errMsg error

func (m *model) startGeneration() {
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
