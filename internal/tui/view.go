package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	if m.isGenerating {
		return m.generationView()
	} else if m.isFinished {
		return m.finishedView()
	}
	return m.configView()
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
	b.WriteString(helpStyle.Render("(esc to quit | tab to navigate)"))

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (m Model) generationView() string {
	var b strings.Builder
	// Use spinner.View() directly without newline to preserve animation
	b.WriteString(fmt.Sprintf("%s Generating addresses...\n\n", m.spinner.View()))
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

func (m Model) finishedView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Attempts: %d\n", m.stats.Attempts))
	b.WriteString(fmt.Sprintf("Found: %d\n", m.stats.Found))

	if len(m.lastResults) > 0 {
		b.WriteString("\nLast found addresses:\n")
		for _, res := range m.lastResults {
			b.WriteString(fmt.Sprintf("• %s\n", res))
		}
	}

	// b.WriteString(helpStyle.Render("\nPress Ctrl+C to exit"))

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(b.String())
}

func NewProgram(m Model) *tea.Program {
	return tea.NewProgram(m)
}
