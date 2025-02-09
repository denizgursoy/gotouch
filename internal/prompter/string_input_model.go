package prompter

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	textInput *textinput.Model
	direction string
	err       error
}

func initialModel(direction, defaultValue string, validator textinput.ValidateFunc) *model {
	ti := textinput.New()
	ti.SetValue(defaultValue)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Validate = validator

	return &model{
		textInput: &ti,
		err:       nil,
		direction: direction,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}
	textInput, cmd := m.textInput.Update(msg)
	m.textInput = &textInput
	return m, cmd
}

func (m *model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.direction,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
