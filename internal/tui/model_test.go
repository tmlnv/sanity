package tui

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
)

func TestInitialModel(t *testing.T) {
	tests := []struct {
		name string
		want Model
	}{
		{
			name: "Initial model with default values",
			want: Model{
				config: config.Config{
					NumAddresses:    1,
					Concurrency:     runtime.NumCPU(),
					LogFile:         config.LogFile,
					PrivateKeysFile: config.PrivateKeysFile,
				},
				inputs: func() []textinput.Model {
					inputs := make([]textinput.Model, 6)
					prompts := []string{
						"Vanity prefix (e.g. '111'):",
						"Vanity suffix:",
						"Vanity regexp:",
						"Number of addresses to find (0=infinite) (default 1):",
						fmt.Sprintf("Workers (0=auto, CPUs available: %d):", runtime.NumCPU()),
						"Timeout (e.g. 30s, 5m):",
					}
					for i := range inputs {
						ti := textinput.New()
						ti.Prompt = prompts[i]
						ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
						inputs[i] = ti
					}
					inputs[0].Focus()
					return inputs
				}(),
				spinner: spinner.New(
					spinner.WithSpinner(spinner.Dot),
					spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("206"))),
				),
				updateChan: make(chan generator.StatsUpdate, 100),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitialModel()
			if !reflect.DeepEqual(got.config, tt.want.config) {
				t.Errorf("InitialModel().config = %v, want %v", got.config, tt.want.config)
			}
			if got.state != tt.want.state {
				t.Errorf("InitialModel().state = %v, want %v", got.state, tt.want.state)
			}
			if len(got.inputs) != len(tt.want.inputs) {
				t.Errorf("InitialModel().inputs length = %v, want %v", len(got.inputs), len(tt.want.inputs))
			}
		})
	}
}

func TestModel_startGeneration(t *testing.T) {
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
		check  func(*testing.T, *Model)
	}{
		{
			name: "Start generation changes state and sets up context",
			fields: fields{
				state: isConfig,
				config: config.Config{
					NumAddresses:    1,
					Concurrency:     runtime.NumCPU(),
					LogFile:         config.LogFile,
					PrivateKeysFile: config.PrivateKeysFile,
				},
				updateChan: make(chan generator.StatsUpdate, 100),
			},
			check: func(t *testing.T, m *Model) {
				if m.state != isGenerating {
					t.Errorf("startGeneration() state = %v, want %v", m.state, isGenerating)
				}
				if m.cancel == nil {
					t.Error("startGeneration() cancel function not set")
				}

				// Test context cancellation
				m.cancel()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
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
			m.startGeneration()
			if tt.check != nil {
				tt.check(t, m)
			}
		})
	}
}
