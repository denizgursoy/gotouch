package prompter

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type textAreaModel struct {
	textarea  *textarea.Model
	err       error
	direction string
}

func initialTextAreaModel(direction, defaultValue string) *textAreaModel {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.SetValue(defaultValue)
	ti.Focus()

	return &textAreaModel{
		textarea:  &ti,
		direction: direction,
		err:       nil,
	}
}

func (m *textAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *textAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	udaptedTextarea, cmd := m.textarea.Update(msg)
	m.textarea = &udaptedTextarea
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *textAreaModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.direction,
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}
