package tui

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/tmlnv/sanity/internal/config"
)

type Model struct {
	step   int
	config config.Config
	inputs []textinput.Model
	err    error
}

func InitialModel() Model {
	m := Model{
		config: config.Config{
			NumAddresses: 1,
			Concurrency:  runtime.NumCPU(),
		},
		inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32
		t.PromptStyle = t.PromptStyle.Faint(true)

		switch i {
		case 0:
			t.Prompt = "Enter vanity prefix: "
			t.Placeholder = "sol"
			t.Validate = validateHexString
		case 1:
			t.Prompt = "Number of addresses to find (0=infinite): "
			t.Validate = validateNumber
		case 2:
			t.Prompt = fmt.Sprintf("Threads to use (%d available): ", runtime.NumCPU())
			t.Validate = validateNumber
		case 3:
			t.Prompt = "Timeout (e.g. 30s, 5m): "
			t.Validate = validateDuration
		}
		m.inputs[i] = t
	}
	return m
}

func validateHexString(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.ParseUint(s, 16, 64)
	return err
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
	return err
}
