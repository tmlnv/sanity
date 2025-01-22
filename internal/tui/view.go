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
)

func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.step >= len(m.inputs) {
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

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Solana Vanity Address Generator"),
		m.inputs[m.step].View(),
		"(esc to quit | enter to continue)",
	)
}

func NewProgram(m Model) *tea.Program {
	return tea.NewProgram(m)
}
