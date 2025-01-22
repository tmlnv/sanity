package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
)

type model struct {
	spinner  spinner.Model
	progress progress.Model
	stats    struct {
		seconds int
		matches int
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := progress.New(progress.WithDefaultGradient())

	return model{
		spinner:  s,
		progress: p,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m.stats.seconds++
		return m, tick()
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	style := lipgloss.NewStyle().Padding(1, 2)
	return style.Render(
		fmt.Sprintf("%s Generating vanity addresses...\nSeconds: %d\nMatches: %d",
			m.spinner.View(), m.stats.seconds, m.stats.matches),
	)
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// StartTUI starts the BubbleTea TUI.
func StartTUI(cfg *config.Config) {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatal("Failed to start TUI", "error", err)
	}
}
