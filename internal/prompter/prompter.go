//go:generate mockgen -source=./prompter.go -destination=mockPrompter.go -package=prompter

package prompter

import (
	"errors"
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	once      = sync.Once{}
	prompter  Prompter
	EmptyList = errors.New("options can not be empty")
)

type (
	Prompter interface {
		AskForString(direction string, validator Validator) (string, error)
		AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error)
		AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error)
		AskForYesOrNo(direction string) (bool, error)
		AskForMultilineString(direction, defaultValue, pattern string) (string, error)
	}

	ListOption struct {
		DisplayText string
		ReturnVal   any
	}

	Validator func(any) error
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = &terminalPrompter{}
	})
	return prompter
}

type terminalPrompter struct {
}

func (tp *terminalPrompter) AskForString(direction string, validator Validator) (string, error) {
	model := initialModel(direction, func(s string) error {
		return validator(s)
	})
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return "nil", fmt.Errorf("error running prompter: %w", err)
	}

	return model.textInput.Value(), nil

}

func (tp *terminalPrompter) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	model := newListModel(direction, list)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return nil, fmt.Errorf("error running prompter: %w", err)
	}
	i, ok := model.list.SelectedItem().(item)
	if !ok {
		return nil, EmptyList
	}

	return i.userData, nil
}

func (tp *terminalPrompter) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	panic("")
}

func (tp *terminalPrompter) AskForYesOrNo(direction string) (bool, error) {
	panic("")
}

func (tp *terminalPrompter) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	panic("")
}
