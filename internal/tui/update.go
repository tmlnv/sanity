package tui

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmlnv/sanity/internal/generator"
)

const (
	prefixStep = iota
	numbAddressesStep
	numThreadsStep
	timeoutStep
	submitStep
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	var r func(tea.Msg) (tea.Model, tea.Cmd)
	switch m.state {
	case isConfig:
		r = m.updateInputs
	case isGenerating:
		r = m.updateGeneration
	case isFinished:
		r = m.updateFinished
	}
	return r(msg)
}

// Add this helper method to handle final updates
func (m *Model) handleStatsUpdate(update generator.StatsUpdate) tea.Cmd {
	m.stats = update.Stats
	if update.LastGenerated != "" {
		m.lastGenerated = append(m.lastGenerated, update.LastGenerated)
		if len(m.lastGenerated) > 5 {
			m.lastGenerated = m.lastGenerated[1:]
		}
	}

	if update.LastMatch != "" {
		m.matched = append(m.matched, update.LastMatch)
	}

	if update.IsFinished {
		m.state = isFinished
		if m.cancel != nil {
			m.cancel() // Properly cancel context
		}
		return tea.Quit
	}

	return m.listenForUpdates()
}

// Update the updateGeneration case
func (m Model) updateGeneration(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case generator.StatsUpdate:
		return m, m.handleStatsUpdate(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateFinished(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case generator.StatsUpdate:
		return m, m.handleStatsUpdate(msg)
	}
	return m, nil
}

func (m *Model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.err = nil

		switch msg.Type {
		case tea.KeyEnter:
			switch m.step {
			case timeoutStep:
				m.step++
			case submitStep:
				m.parseInputs()
				go m.startGeneration()
				return m, tea.Batch(
					m.spinner.Tick,
					m.listenForUpdates(),
				)
			}
			return m.handleInput()

		case tea.KeyRunes:
			switch m.step {
			case timeoutStep:
				return m.handleDurationInput(msg)
			}

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown:
			return m.handleNavigation(msg)
		}
	}

	return m.updateFocusedInput(msg)
}

func (m *Model) handleDurationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	current := m.inputs[timeoutStep].Value()

	// Allow numeric input at any position
	if _, err := strconv.Atoi(string(msg.Runes)); err == nil {
		return m.updateFocusedInput(msg)
	}

	// Allow unit characters only at end
	if strings.ContainsAny(string(msg.Runes), "smh") {
		if len(current) == 0 || strings.ContainsAny(current[len(current)-1:], "0123456789") {
			return m.updateFocusedInput(msg)
		}
	}

	// Block invalid characters
	m.err = fmt.Errorf("only numbers and s/m/h units are allowed in duration")
	return m, nil
}

func (m *Model) parseInputs() {
	m.config.Prefix = m.inputs[0].Value()

	if val := m.inputs[1].Value(); val != "" {
		m.config.NumAddresses, _ = strconv.Atoi(val)
	}

	if val := m.inputs[2].Value(); val != "" {
		if threads, _ := strconv.Atoi(val); threads > 0 {
			m.config.Concurrency = threads
		} else {
			m.config.Concurrency = runtime.NumCPU()
		}
	}

	if val := m.inputs[3].Value(); val != "" {
		// Handle pure numeric input as seconds
		if dur, err := strconv.Atoi(val); err == nil {
			m.config.Timeout = time.Duration(dur) * time.Second
		} else {
			m.config.Timeout, _ = time.ParseDuration(val)
		}
	}
}

func (m *Model) handleInput() (*Model, tea.Cmd) {
	if m.step < len(m.inputs)-1 {
		m.step++
		m.inputs[m.step].Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m *Model) handleNavigation(msg tea.KeyMsg) (*Model, tea.Cmd) {
	s := msg.String()
	if s == "up" || s == "shift+tab" {
		m.step--
	} else {
		m.step++
	}

	// Keep step in valid range [0, submitStep]
	if m.step > submitStep {
		m.step = 0
	} else if m.step < 0 {
		m.step = submitStep
	}

	// Only update focus for input fields
	if m.step < len(m.inputs) {
		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i < len(m.inputs); i++ {
			if i == m.step {
				cmds[i] = m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *Model) updateFocusedInput(msg tea.Msg) (*Model, tea.Cmd) {
	if m.step >= len(m.inputs) {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
	return m, cmd
}

func (m Model) listenForUpdates() tea.Cmd {
	return func() tea.Msg {
		update, ok := <-m.updateChan
		if !ok {
			return tea.Quit
		}
		return update
	}
}
