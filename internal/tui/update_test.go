package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmlnv/sanity/internal/generator"
)

func TestModel_Init(t *testing.T) {
	m := Model{}
	got := m.Init()
	if got == nil {
		t.Errorf("Model.Init() = %v, want non-nil tea.Cmd", got)
	}
}

func TestModel_Update_KeyMsg(t *testing.T) {
	m := Model{state: isConfig}
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	got, cmd := m.Update(msg)
	// Can't compare structs with slices directly, so just check type
	if _, ok := got.(Model); !ok {
		t.Errorf("Model.Update() got = %T, want Model", got)
	}
	if cmd == nil {
		t.Errorf("Model.Update() cmd = %v, want non-nil", cmd)
	}
}

func TestModel_handleStatsUpdate_Finished(t *testing.T) {
	called := false
	cancel := func() { called = true }
	m := &Model{
		cancel: cancel,
	}
	update := generator.StatsUpdate{
		Stats:      generator.Stats{Attempts: 1, Found: 1},
		LastMatch:  "match",
		IsFinished: true,
	}
	cmd := m.handleStatsUpdate(update)
	if m.state != isFinished {
		t.Errorf("handleStatsUpdate did not set state to isFinished")
	}
	if !called {
		t.Errorf("handleStatsUpdate did not call cancel")
	}
	if cmd == nil {
		t.Errorf("handleStatsUpdate did not return tea.Quit")
	}
}

func TestModel_updateGeneration_StatsUpdate(t *testing.T) {
	m := Model{}
	msg := generator.StatsUpdate{Stats: generator.Stats{Attempts: 1}}
	got, cmd := m.updateGeneration(msg)
	// Can't compare structs with slices directly, so just check type
	if _, ok := got.(Model); !ok {
		t.Errorf("updateGeneration got = %T, want Model", got)
	}
	if cmd == nil {
		t.Errorf("updateGeneration cmd = %v, want non-nil", cmd)
	}
}

func TestModel_updateGeneration_SpinnerTickMsg(t *testing.T) {
	m := Model{spinner: spinner.New()}
	msg := spinner.TickMsg{}
	got, cmd := m.updateGeneration(msg)
	// Can't compare structs with slices directly, so just check type
	if _, ok := got.(Model); !ok {
		t.Errorf("updateGeneration got = %T, want Model", got)
	}
	if cmd == nil {
		t.Errorf("updateGeneration cmd = %v, want non-nil", cmd)
	}
}

func TestModel_updateFinished_StatsUpdate(t *testing.T) {
	m := Model{}
	msg := generator.StatsUpdate{Stats: generator.Stats{Attempts: 1}}
	got, cmd := m.updateFinished(msg)
	// Can't compare structs with slices directly, so just check type
	if _, ok := got.(Model); !ok {
		t.Errorf("updateFinished got = %T, want Model", got)
	}
	if cmd == nil {
		t.Errorf("updateFinished cmd = %v, want non-nil", cmd)
	}
}

func TestModel_updateInputs_Enter(t *testing.T) {
	m := &Model{
		step:   timeoutStep,
		inputs: make([]textinput.Model, submitStep+1),
	}
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	got, _ := m.updateInputs(msg)
	if got == nil {
		t.Errorf("updateInputs should not return nil")
	}
}

func TestModel_handleDurationInput_Number(t *testing.T) {
	m := &Model{
		inputs: make([]textinput.Model, timeoutStep+1),
	}
	m.inputs[timeoutStep] = textinput.New()
	msg := tea.KeyMsg{Runes: []rune("1")}
	got, _ := m.handleDurationInput(msg)
	if got == nil {
		t.Errorf("handleDurationInput should not return nil")
	}
}

func TestModel_parseInputs_Valid(t *testing.T) {
	m := &Model{
		inputs: make([]textinput.Model, timeoutStep+1),
	}
	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}
	m.inputs[numbAddressesStep].SetValue("2")
	m.inputs[numThreadsStep].SetValue("2")
	m.inputs[timeoutStep].SetValue("1")
	err := m.parseInputs()
	if err != nil {
		t.Errorf("parseInputs() error = %v, want nil", err)
	}
	if m.config.NumAddresses != 2 {
		t.Errorf("parseInputs() NumAddresses = %v, want 2", m.config.NumAddresses)
	}
	if m.config.Concurrency != 2 {
		t.Errorf("parseInputs() Concurrency = %v, want 2", m.config.Concurrency)
	}
	if m.config.Timeout != 1 {
		t.Errorf("parseInputs() Timeout = %v, want 1", m.config.Timeout)
	}
}

func TestModel_handleInput_StepAdvance(t *testing.T) {
	m := &Model{
		inputs: make([]textinput.Model, submitStep+1),
		step:   0,
	}
	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}
	got, _ := m.handleInput()
	if got.step != 1 {
		t.Errorf("handleInput() step = %v, want 1", got.step)
	}
}

func TestModel_handleNavigation_UpDown(t *testing.T) {
	m := &Model{
		inputs: make([]textinput.Model, submitStep+1),
		step:   0,
	}
	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}
	msg := tea.KeyMsg{Type: tea.KeyUp}
	got, _ := m.handleNavigation(msg)
	if got == nil {
		t.Errorf("handleNavigation() should not return nil")
	}
}

func TestModel_updateFocusedInput(t *testing.T) {
	m := &Model{
		inputs: make([]textinput.Model, submitStep+1),
		step:   0,
	}
	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}
	got, _ := m.updateFocusedInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if got == nil {
		t.Errorf("updateFocusedInput() should not return nil")
	}
}

func TestModel_listenForUpdates(t *testing.T) {
	m := Model{
		updateChan: make(chan generator.StatsUpdate, 1),
	}
	m.updateChan <- generator.StatsUpdate{Stats: generator.Stats{Attempts: 1}}
	cmd := m.listenForUpdates()
	if cmd == nil {
		t.Errorf("listenForUpdates() should not return nil")
	}
}
