package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.generating {
		return m.generationView()
	}

	if m.step >= len(m.inputs) {
		return m.configCompleteView()
	}

	return m.configView()
}

func (m Model) configView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Solana Vanity Address Generator"),
		m.inputs[m.step].View(),
		helpStyle.Render("(esc to quit | enter to continue)"),
	)
}

func (m Model) configCompleteView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Configuration Complete!"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Prefix: %s\n", m.config.Prefix))
	b.WriteString(fmt.Sprintf("Count: %d\n", m.config.NumAddresses))
	b.WriteString(fmt.Sprintf("Threads: %d\n", m.config.Concurrency))
	b.WriteString(fmt.Sprintf("Timeout: %s\n", m.config.Timeout))
	b.WriteString("\nPress any key to start generation...")
	return b.String()
}

func (m Model) generationView() string {
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

func NewProgram(m Model) *tea.Program {
	return tea.NewProgram(m)
}
