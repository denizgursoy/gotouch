package prompter

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/manifoldco/promptui"
	"os"
	"syscall"
)

type promptUi struct {
	Manager manager.Manager
}

func (p *promptUi) AskForSelectionFromList(direction string, list []fmt.Stringer) (interface{}, error) {
	if !isValid(p) {
		return "", model.ErrMissingField
	}

	count := len(list)

	if count == 0 {
		return nil, EmptyList
	} else if count == 1 {
		return list[0], nil
	}

	options := make([]string, 0)
	for _, item := range list {
		options = append(options, item.String())
	}

	prompt := promptui.Select{
		Label: direction,
		Items: options,
		Stdin: p.Manager.GetStream(),
	}

	index, _, err := prompt.Run()
	if err != nil {
		p.exitIfInterrupted(err)
		return nil, err
	}

	return list[index], nil
}

func (p *promptUi) exitIfInterrupted(err error) {
	if err.Error() == "^C" || err.Error() == "^D" {
		os.Exit(int(syscall.SIGINT))
	}
}

func (p *promptUi) AskForString(direction string, validator StringValidator) (string, error) {
	if !isValid(p) {
		return "", model.ErrMissingField
	}

	prompt := promptui.Prompt{
		Label:    direction,
		Validate: promptui.ValidateFunc(validator),
		Stdin:    p.Manager.GetStream(),
	}
	input, err := prompt.Run()
	if err != nil {
		p.exitIfInterrupted(err)
		return "", err
	}
	return input, nil
}

func isValid(p *promptUi) bool {
	return p.Manager != nil
}
