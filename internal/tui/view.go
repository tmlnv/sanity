package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle    = blurredStyle.Copy()
	titleStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	mainStyle  = lipgloss.NewStyle()
)

func (m Model) View() string {
	var s string
	if m.err != nil {
		s = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	switch m.state {
	case isConfig:
		s = m.configView()
	case isGenerating:
		s = m.generationView()
	case isFinished:
		s = m.finishedView()
	}
	return s
}

func (m Model) configView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Solana Vanity Address Generator"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredStyle.Render("[ Submit ]")
	if m.step == submitStep {
		button = focusedStyle.Render("[ Submit ]")
	}
	b.WriteString("\n\n" + button + "\n")

	b.WriteString("\n")
	if m.err != nil {
		b.WriteString(helpStyle.Render("Error: ") + errorStyle.Render(m.err.Error()) + "\n")
	}

	b.WriteString(helpStyle.Render("(esc to quit | tab to navigate)"))

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (m Model) generationView() string {
	var b strings.Builder
	// Using spinner.View() directly without newline to preserve animation
	b.WriteString(mainStyle.Render(m.spinner.View() + "Generating addresses...\n\n"))
	b.WriteString(mainStyle.Render(m.config.String()) + "\n")

	if len(m.lastGenerated) > 0 {
		for _, addr := range m.lastGenerated {
			b.WriteString(mainStyle.Render(fmt.Sprintf("• %s\n", addr)))
		}
	}
	b.WriteString("\n")

	b.WriteString(mainStyle.Render("Attempts: " + strconv.Itoa(int(m.stats.Attempts)) + "\n"))
	b.WriteString(mainStyle.Render(fmt.Sprintf("Found: %d\n", m.stats.Found)))
	if len(m.matched) > 0 {
		b.WriteString("\nFound addresses:\n")
		for _, addr := range m.matched {
			b.WriteString(mainStyle.Render(fmt.Sprintf("• %s\n", addr)))
		}
	}

	b.WriteString(helpStyle.Render("\nPress Esc or Ctrl+C to exit"))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(b.String())
}

func (m Model) finishedView() string {
	var b strings.Builder
	b.WriteString("Finished\n\n")
	b.WriteString(fmt.Sprintf("Attempts: %d\n", m.stats.Attempts))
	b.WriteString(fmt.Sprintf("Found: %d\n", m.stats.Found))

	if len(m.matched) > 0 {
		b.WriteString("\nFound addresses:\n")
		for _, addr := range m.matched {
			b.WriteString(fmt.Sprintf("• %s\n", addr))
		}
	}

	b.WriteString(fmt.Sprintf("\nCorresponding private keys were written to %v\n", config.PrivateKeysFile))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(b.String())
}

func NewProgram(m Model) *tea.Program {
	return tea.NewProgram(m)
}
