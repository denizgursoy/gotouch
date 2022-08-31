package prompter

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/denizgursoy/gotouch/internal/manager"
)

type (
	srv struct {
		Manager manager.Manager
	}
)

func (s srv) AskForString(direction string, validator StringValidator) (string, error) {
	result := ""
	input := survey.Input{
		Message: direction,
	}
	err := survey.AskOne(&input, &result)
	return result, err
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (interface{}, error) {
	//if !isValid(s) {
	//	return "", model.ErrMissingField
	//}

	count := len(list)

	if count == 0 {
		return nil, EmptyList
	} else if count == 1 {
		return list[0], nil
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

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	name := false
	prompt := &survey.Confirm{
		Message: direction,
	}
	err := survey.AskOne(prompt, &name)
	return name, err
}
