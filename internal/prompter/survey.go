//go:build !integration

package prompter

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func (s srv) AskForString(direction string, validator Validator) (string, error) {
	result := ""

	input := survey.Input{
		Message: direction,
	}

	err := survey.AskOne(&input, &result, survey.WithValidator(survey.Validator(validator)))
	return result, err
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make(map[string]fmt.Stringer)
	keys := make([]string, 0)
	for _, item := range list {
		choice := item.String()
		options[choice] = item
		keys = append(keys, choice)
	}

	selectedChoice := ""
	err := survey.AskOne(&survey.Select{
		Message: direction,
		Options: keys,
	}, &selectedChoice)

	return options[selectedChoice], err
}

func (s srv) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	count := len(list)

	if count == 0 {
		return nil, EmptyList
	}

	options := make(map[string]fmt.Stringer)
	keys := make([]string, 0)
	for _, item := range list {
		choice := item.String()
		options[choice] = item
		keys = append(keys, choice)
	}

	selectedChoices := make([]string, 0)
	err := survey.AskOne(&survey.MultiSelect{
		Message: direction,
		Options: keys,
	}, &selectedChoices)

	results := make([]any, 0)
	for i := range selectedChoices {
		results = append(results, options[selectedChoices[i]])
	}

	return results, err
}

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: direction,
	}
	err := survey.AskOne(prompt, &name)
	return name, err
}

func (s srv) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	prompt := &survey.Editor{
		Message:       direction,
		Default:       defaultValue,
		HideDefault:   true,
		AppendDefault: true,
		FileName:      pattern,
	}

	result := ""
	err := survey.AskOne(prompt, &result)
	return result, err
}
