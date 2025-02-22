package tui

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
)

func TestModel_View(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		check  func(t *testing.T, got string)
	}{
		{
			name: "error view",
			fields: fields{
				state: isConfig,
				inputs: []textinput.Model{
					textinput.New(),
				},
				err: fmt.Errorf("test error"),
			},
			check: func(t *testing.T, got string) {
				// Check for key elements in the output
				if !strings.Contains(got, "Solana Vanity Address Generator") {
					t.Error("Missing title")
				}
				if !strings.Contains(got, "Error: test error") {
					t.Error("Missing error message")
				}
				if !strings.Contains(got, "[ Submit ]") {
					t.Error("Missing submit button")
				}
				if !strings.Contains(got, "(esc to quit | tab to navigate)") {
					t.Error("Missing help text")
				}
			},
		},
		{
			name: "config state",
			fields: fields{
				state: isConfig,
				inputs: []textinput.Model{
					textinput.New(),
				},
				step: submitStep,
			},
			check: func(t *testing.T, got string) {
				want := lipgloss.NewStyle().Padding(1, 2).Render(
					titleStyle.Render("Solana Vanity Address Generator") + "\n\n" +
						">" + "\n\n" +
						focusedStyle.Render("[ Submit ]") + "\n\n" +
						helpStyle.Render("(esc to quit | tab to navigate)"))
				if got != want {
					t.Errorf("got = %v, want %v", got, want)
				}
			},
		},
		{
			name: "generating state",
			fields: fields{
				state:   isGenerating,
				spinner: spinner.Model{},
				config:  config.Config{},
				stats: generator.Stats{
					Attempts: 100,
					Found:    2,
				},
				lastGenerated: []string{"addr1", "addr2"},
				matched:       []string{"matched1"},
			},
			check: func(t *testing.T, got string) {
				want := lipgloss.NewStyle().Padding(1, 2).Render(
					mainStyle.Render("(error)Generating addresses...\n\n") +
						mainStyle.Render("Prefix:\nSuffix:\nRegexp:\nNumAddresses: 0\nConcurrency: 0\nTimeout: 0s\nLogFile:\nPrivateKeysFile:\n\n") +
						mainStyle.Render("• addr1\n") +
						mainStyle.Render("• addr2\n") +
						"\n" +
						mainStyle.Render("Attempts: 100\n") +
						mainStyle.Render("Found: 2\n") +
						"\nFound addresses:\n" +
						mainStyle.Render("• matched1\n") +
						helpStyle.Render("\nPress Esc or Ctrl+C to exit"))
				if got != want {
					t.Errorf("got = %v, want %v", got, want)
				}
			},
		},
		{
			name: "finished state",
			fields: fields{
				state: isFinished,
				stats: generator.Stats{
					Attempts: 200,
					Found:    3,
				},
				matched: []string{"addr1", "addr2", "addr3"},
			},
			check: func(t *testing.T, got string) {
				want := lipgloss.NewStyle().Padding(1, 2).Render(
					"Finished\n\n" +
						"Attempts: 200\n" +
						"Found: 3\n\n" +
						"Found addresses:\n" +
						"• addr1\n" +
						"• addr2\n" +
						"• addr3\n\n" +
						fmt.Sprintf("Corresponding private keys were written to %v:\n", config.PrivateKeysFile))
				if got != want {
					t.Errorf("got = %v, want %v", got, want)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				inputs:        tt.fields.inputs,
				step:          tt.fields.step,
				spinner:       tt.fields.spinner,
				state:         tt.fields.state,
				config:        tt.fields.config,
				stats:         tt.fields.stats,
				lastGenerated: tt.fields.lastGenerated,
				matched:       tt.fields.matched,
				err:           tt.fields.err,
				updateChan:    tt.fields.updateChan,
				cancel:        tt.fields.cancel,
			}
			tt.check(t, m.View())
		})
	}
}

func TestModel_configView(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "config view with no error",
			fields: fields{
				inputs: []textinput.Model{
					textinput.New(),
				},
				step: submitStep,
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				titleStyle.Render("Solana Vanity Address Generator") + "\n\n" +
					">" + "\n\n" +
					focusedStyle.Render("[ Submit ]") + "\n\n" +
					helpStyle.Render("(esc to quit | tab to navigate)")),
		},
		{
			name: "config view with error",
			fields: fields{
				inputs: []textinput.Model{
					textinput.New(),
				},
				step: submitStep,
				err:  fmt.Errorf("test error"),
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				titleStyle.Render("Solana Vanity Address Generator") + "\n\n" +
					">" + "\n\n" +
					focusedStyle.Render("[ Submit ]") + "\n\n" +
					helpStyle.Render("Error: ") + errorStyle.Render("test error") + "\n" +
					helpStyle.Render("(esc to quit | tab to navigate)")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				inputs:        tt.fields.inputs,
				step:          tt.fields.step,
				spinner:       tt.fields.spinner,
				state:         tt.fields.state,
				config:        tt.fields.config,
				stats:         tt.fields.stats,
				lastGenerated: tt.fields.lastGenerated,
				matched:       tt.fields.matched,
				err:           tt.fields.err,
				updateChan:    tt.fields.updateChan,
				cancel:        tt.fields.cancel,
			}
			if got := m.configView(); got != tt.want {
				t.Errorf("Model.configView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_generationView(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "generation view with no matches",
			fields: fields{
				spinner: spinner.Model{},
				config:  config.Config{},
				stats: generator.Stats{
					Attempts: 50,
					Found:    0,
				},
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				mainStyle.Render("(error)Generating addresses...\n\n") +
					mainStyle.Render("Prefix:\nSuffix:\nRegexp:\nNumAddresses: 0\nConcurrency: 0\nTimeout: 0s\nLogFile:\nPrivateKeysFile:\n\n") +
					"\n" +
					mainStyle.Render("Attempts: 50\n") +
					mainStyle.Render("Found: 0\n") +
					helpStyle.Render("\nPress Esc or Ctrl+C to exit")),
		},
		{
			name: "generation view with matches",
			fields: fields{
				spinner: spinner.Model{},
				config:  config.Config{},
				stats: generator.Stats{
					Attempts: 100,
					Found:    2,
				},
				lastGenerated: []string{"addr1", "addr2"},
				matched:       []string{"matched1", "matched2"},
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				mainStyle.Render("(error)Generating addresses...\n\n") +
					mainStyle.Render("Prefix:\nSuffix:\nRegexp:\nNumAddresses: 0\nConcurrency: 0\nTimeout: 0s\nLogFile:\nPrivateKeysFile:\n\n") +
					mainStyle.Render("• addr1\n") +
					mainStyle.Render("• addr2\n") +
					"\n" +
					mainStyle.Render("Attempts: 100\n") +
					mainStyle.Render("Found: 2\n") +
					"\nFound addresses:\n" +
					mainStyle.Render("• matched1\n") +
					mainStyle.Render("• matched2\n") +
					helpStyle.Render("\nPress Esc or Ctrl+C to exit")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				inputs:        tt.fields.inputs,
				step:          tt.fields.step,
				spinner:       tt.fields.spinner,
				state:         tt.fields.state,
				config:        tt.fields.config,
				stats:         tt.fields.stats,
				lastGenerated: tt.fields.lastGenerated,
				matched:       tt.fields.matched,
				err:           tt.fields.err,
				updateChan:    tt.fields.updateChan,
				cancel:        tt.fields.cancel,
			}
			if got := m.generationView(); got != tt.want {
				t.Errorf("Model.generationView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_finishedView(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "finished view with no matches",
			fields: fields{
				stats: generator.Stats{
					Attempts: 100,
					Found:    0,
				},
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				"Finished\n\n" +
					"Attempts: 100\n" +
					"Found: 0\n\n" +
					fmt.Sprintf("Corresponding private keys were written to %v:\n", config.PrivateKeysFile)),
		},
		{
			name: "finished view with matches",
			fields: fields{
				stats: generator.Stats{
					Attempts: 200,
					Found:    3,
				},
				matched: []string{"addr1", "addr2", "addr3"},
			},
			want: lipgloss.NewStyle().Padding(1, 2).Render(
				"Finished\n\n" +
					"Attempts: 200\n" +
					"Found: 3\n\n" +
					"Found addresses:\n" +
					"• addr1\n" +
					"• addr2\n" +
					"• addr3\n\n" +
					fmt.Sprintf("Corresponding private keys were written to %v:\n", config.PrivateKeysFile)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				inputs:        tt.fields.inputs,
				step:          tt.fields.step,
				spinner:       tt.fields.spinner,
				state:         tt.fields.state,
				config:        tt.fields.config,
				stats:         tt.fields.stats,
				lastGenerated: tt.fields.lastGenerated,
				matched:       tt.fields.matched,
				err:           tt.fields.err,
				updateChan:    tt.fields.updateChan,
				cancel:        tt.fields.cancel,
			}
			if got := m.finishedView(); got != tt.want {
				t.Errorf("Model.finishedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProgram(t *testing.T) {
	m := Model{}
	got := NewProgram(m)

	if got == nil {
		t.Error("NewProgram() returned nil")
		return
	}

	// We can't directly compare tea.Program instances
	// Just verify that NewProgram returns a non-nil program
	if reflect.TypeOf(got) != reflect.TypeOf(&tea.Program{}) {
		t.Errorf("NewProgram() returned %T, want *tea.Program", got)
	}
}
