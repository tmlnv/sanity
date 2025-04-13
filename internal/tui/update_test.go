package tui

import (
	"context"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/generator"
)

func TestModel_Init(t *testing.T) {
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
		want   tea.Cmd
	}{
		// TODO: Add test cases.
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
			if got := m.Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_Update(t *testing.T) {
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
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.Update(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.Update() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.Update() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_handleStatsUpdate(t *testing.T) {
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
	type args struct {
		update generator.StatsUpdate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Cmd
	}{
		// TODO: Add test cases.
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
			if got := m.handleStatsUpdate(tt.args.update); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleStatsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_updateGeneration(t *testing.T) {
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
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.updateGeneration(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.updateGeneration() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.updateGeneration() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_updateFinished(t *testing.T) {
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
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.updateFinished(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.updateFinished() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.updateFinished() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_updateInputs(t *testing.T) {
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
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.updateInputs(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.updateInputs() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.updateInputs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_handleDurationInput(t *testing.T) {
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
	type args struct {
		msg tea.KeyMsg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tea.Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.handleDurationInput(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleDurationInput() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleDurationInput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_parseInputs(t *testing.T) {
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
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := m.parseInputs(); (err != nil) != tt.wantErr {
				t.Errorf("Model.parseInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModel_handleInput(t *testing.T) {
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
		want   *Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.handleInput()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleInput() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleInput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_handleNavigation(t *testing.T) {
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
	type args struct {
		msg tea.KeyMsg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.handleNavigation(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.handleNavigation() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.handleNavigation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_updateFocusedInput(t *testing.T) {
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
	type args struct {
		msg tea.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Model
		want1  tea.Cmd
	}{
		// TODO: Add test cases.
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
			got, got1 := m.updateFocusedInput(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.updateFocusedInput() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Model.updateFocusedInput() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestModel_listenForUpdates(t *testing.T) {
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
		want   tea.Cmd
	}{
		// TODO: Add test cases.
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
			if got := m.listenForUpdates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.listenForUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}
