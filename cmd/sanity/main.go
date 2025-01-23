package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
	"github.com/tmlnv/sanity/internal/logger"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle    = blurredStyle.Copy()
	version      = "v1.0.0"
)

func main() {
	cfg := parseFlags()

	if !cfg.FlagsProvided {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			logger.Error("Error running TUI", "error", err)
			os.Exit(1)
		}
		return
	}

	runCLI(cfg)
}

func parseFlags() config.Config {
	var cfg config.Config

	flag.StringVar(&cfg.Prefix, "prefix", "", "Vanity prefix for Solana address")
	flag.StringVar(&cfg.Suffix, "suffix", "", "Vanity suffix for Solana address")
	flag.StringVar(&cfg.Regex, "regex", "", "Regex pattern to match")
	flag.IntVar(&cfg.NumAddresses, "count", 1, "Number of addresses to find (0=infinite)")
	flag.IntVar(&cfg.Concurrency, "threads", 0, "Number of worker threads (0=auto)")
	flag.DurationVar(&cfg.Timeout, "timeout", 0, "Maximum search duration")
	flag.StringVar(&cfg.LogFile, "logfile", "", "Path to log file")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("sanity %s\n", version)
		os.Exit(0)
	}

	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.NumCPU()
	}

	cfg.FlagsProvided = flag.NFlag() > 0
	return cfg
}

func runCLI(cfg config.Config) {
	// Initialize logger before starting generator
	logger.Init(cfg.LogFile)

	ctx, cancel := createContext(cfg)
	defer cancel()

	logger.Info("Starting vanity generation",
		"prefix", cfg.Prefix,
		"suffix", cfg.Suffix,
		"regex", cfg.Regex,
		"threads", cfg.Concurrency,
		"timeout", cfg.Timeout,
	)

	// Start generator with a nil update channel (CLI mode)
	generator.Start(ctx, cfg, nil)
}

func createContext(cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Timeout > 0 {
		return context.WithTimeout(context.Background(), cfg.Timeout)
	}
	return context.WithCancel(context.Background())
}

type model struct {
	inputs      []textinput.Model
	focusIndex  int
	spinner     spinner.Model
	generating  bool
	config      config.Config
	stats       generator.Stats
	lastResults []string
	updateChan  chan generator.StatsUpdate
}

func initialModel() model {
	m := model{
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
		ti.PromptStyle = blurredStyle
		if in.validation != nil {
			ti.Validate = in.validation
		}
		m.inputs[in.index] = ti
	}

	m.inputs[0].Focus()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.generating {
		return m.updateGeneration(msg)
	}
	return m.updateInputs(msg)
}

func (m model) updateGeneration(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case generator.StatsUpdate:
		m.stats = msg.Stats
		if msg.LastResult != "" {
			m.lastResults = append(m.lastResults, msg.LastResult)
			if len(m.lastResults) > 5 {
				m.lastResults = m.lastResults[1:]
			}
		}
		return m, m.listenForUpdates()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				m.parseInputs()
				m.generating = true
				go m.startGeneration()
				return m, tea.Batch(
					m.spinner.Tick,
					m.listenForUpdates(),
				)
			}

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown:
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateFocusedInput(msg)
	return m, cmd
}

func (m *model) parseInputs() {
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
		m.config.Timeout, _ = time.ParseDuration(val)
	}
}

func (m *model) updateFocusedInput(msg tea.Msg) tea.Cmd {
	if m.focusIndex >= len(m.inputs) {
		return nil
	}

	input := m.inputs[m.focusIndex]
	var cmd tea.Cmd
	input, cmd = input.Update(msg)
	m.inputs[m.focusIndex] = input
	return cmd
}

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

func (m model) listenForUpdates() tea.Cmd {
	return func() tea.Msg {
		update, ok := <-m.updateChan
		if !ok {
			// Channel closed, generation completed
			return tea.Quit
		}
		return update
	}
}

func (m model) View() string {
	if m.generating {
		return m.generationView()
	}
	return m.configView()
}

func (m model) configView() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredStyle.Render("[ Submit ]")
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
	}
	b.WriteString("\n\n" + button + "\n")

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(b.String())
}

func (m model) generationView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s Generating addresses...\n\n", m.spinner.View()))
	b.WriteString(fmt.Sprintf("Attempts: %d\n", m.stats.Attempts))
	b.WriteString(fmt.Sprintf("Found: %d\n", m.stats.Found))

	if len(m.lastResults) > 0 {
		b.WriteString("\nLast found addresses:\n")
		for _, res := range m.lastResults {
			b.WriteString(fmt.Sprintf("• %s\n", res))
		}
	}

	b.WriteString(helpStyle.Render("\nPress Ctrl+C to exit"))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(b.String())
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
