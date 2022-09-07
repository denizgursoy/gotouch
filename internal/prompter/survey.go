package prompter

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/denizgursoy/gotouch/internal/manager"
)

type (
	srv struct {
		Manager manager.Manager
	}
)

func (s srv) AskForString(direction string, validator Validator) (string, error) {
	if s.Manager.IsTest() {
		all, err := ioutil.ReadAll(s.Manager.GetStream())
		if err != nil {
			return "", err
		}
		return string(all), nil
	}

	result := ""

	input := survey.Input{
		Message: direction,
	}

	err := survey.AskOne(&input, &result, survey.WithValidator(survey.Validator(validator)))
	return result, err
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (interface{}, error) {
	if s.Manager.IsTest() {
		all, err := ioutil.ReadAll(s.Manager.GetStream())
		if err != nil {
			return "", err
		}

		atoi, err := strconv.Atoi(string(all))
		return list[atoi], nil
	}

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
	if s.Manager.IsTest() {
		all, err := ioutil.ReadAll(s.Manager.GetStream())
		if err != nil {
			return false, err
		}
		atoi, err := strconv.Atoi(string(all))

		return atoi == 1, nil
	}

	name := false
	prompt := &survey.Confirm{
		Message: direction,
	}
	err := survey.AskOne(prompt, &name)
	return name, err
}
