//go:build !integration

package prompter

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func newTheme() *huh.Theme {
	theme := huh.ThemeCharm()
	theme.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")
	theme.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	theme.Blurred.SelectedPrefix = lipgloss.NewStyle().SetString("[✓] ")
	theme.Blurred.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	return theme
}

func runInline(fields ...huh.Field) error {
	return huh.NewForm(huh.NewGroup(fields...)).
		WithProgramOptions(tea.WithAltScreen()).
		WithTheme(newTheme()).
		Run()
}

func (s srv) AskForString(direction, initialValue string, validator Validator) (string, error) {
	result := initialValue

	input := huh.NewInput().
		Title(direction).
		Value(&result)

	if validator != nil {
		input.Validate(func(val string) error {
			return validator(val)
		})
	}

	err := runInline(input)
	return result, err
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make([]huh.Option[string], 0, count)
	lookup := make(map[string]fmt.Stringer)
	for _, item := range list {
		choice := item.String()
		options = append(options, huh.NewOption(choice, choice))
		lookup[choice] = item
	}

	var selected string
	err := runInline(
		huh.NewSelect[string]().
			Title(direction + " (select one)").
			Options(options...).
			Value(&selected),
	)

	return lookup[selected], err
}

func (s srv) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make([]huh.Option[string], 0, count)
	lookup := make(map[string]fmt.Stringer)
	for _, item := range list {
		choice := item.String()
		options = append(options, huh.NewOption(choice, choice))
		lookup[choice] = item
	}

	var selected []string
	err := runInline(
		huh.NewMultiSelect[string]().
			Title(direction + " (select multiple)").
			Options(options...).
			Value(&selected),
	)

	results := make([]any, 0, len(selected))
	for _, s := range selected {
		results = append(results, lookup[s])
	}

	return results, err
}

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	var result bool
	err := runInline(
		huh.NewConfirm().
			Title(direction).
			Value(&result),
	)
	return result, err
}

func (s srv) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	result := defaultValue
	err := runInline(
		huh.NewText().
			Title(direction).
			Value(&result),
	)
	return result, err
}
