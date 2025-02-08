package prompter

import tea "github.com/charmbracelet/bubbletea"

type StringInputModel struct {
	direction string
	validator Validator
}

func NewStringInputModel(direction string, validator Validator) *StringInputModel {
	return &StringInputModel{direction: direction, validator: validator}
}

func (m StringInputModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m StringInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			return m, tea.ClearScreen
		}
	}

	// Return the updated bubbleTeaContext to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m StringInputModel) View() string {
	return m.direction
}
