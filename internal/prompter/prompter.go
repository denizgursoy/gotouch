//go:generate mockgen -source=./prompter.go -destination=mockPrompter.go -package=prompter

package prompter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	once      = sync.Once{}
	prompter  Prompter
	EmptyList = errors.New("options can not be empty")
)

type (
	Prompter interface {
		AskForString(direction, defaultValue string, validator Validator) (string, error)
		AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error)
		AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error)
		AskForYesOrNo(direction string) (bool, error)
		AskForMultilineString(direction, defaultValue, pattern string) (string, error)
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

func (tp *terminalPrompter) AskForString(direction, defaultValue string, validator Validator) (string, error) {
	model := initialModel(direction, defaultValue, func(s string) error {
		return validator(s)
	})
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return "nil", fmt.Errorf("error running prompter: %w", err)
	}

	return model.textInput.Value(), nil

}

func (tp *terminalPrompter) AskForSelectionFromList(direction string, itemsToSelect []fmt.Stringer) (any, error) {

	items := make([]list.Item, 0)
	for _, selection := range itemsToSelect {
		items = append(items, item{
			title:       selection.String(),
			description: "test",
			userData:    selection,
		})
	}

	model := newListModel(direction, items)
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
	items := []list.Item{
		item{
			title:       "Yes",
			description: "Absolutely!",
			userData:    true,
		},
		item{
			title:       "No",
			description: "Of course not",
			userData:    false,
		},
	}

	model := newListModel(direction, items)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return false, fmt.Errorf("error running prompter: %w", err)
	}

	i, ok := model.list.SelectedItem().(item)
	if !ok {
		return false, EmptyList
	}

	return i.userData.(bool), nil
}

func (tp *terminalPrompter) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	model := initialTextAreaModel(direction, defaultValue)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return "nil", fmt.Errorf("error running prompter: %w", err)
	}

	return model.textarea.Value(), nil
}
