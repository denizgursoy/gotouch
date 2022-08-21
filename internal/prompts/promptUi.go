package prompts

import (
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/manifoldco/promptui"
	"log"
)

type promptUi struct {
}

func (p promptUi) AskForSelectionFromList(direction string, listOptions []*ListOption) interface{} {
	options := make([]string, 0)
	for _, option := range listOptions {
		options = append(options, option.DisplayText)
	}

	prompt := promptui.Select{
		Label: direction,
		Items: options,
		Stdin: manager.GetInstance().GetStream(),
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Println(err)
	}

	return listOptions[index].ReturnVal
}

func (p promptUi) AskForString(direction string, validator StringValidator) string {
	prompt := promptui.Prompt{
		Label:    direction,
		Validate: promptui.ValidateFunc(validator),
		Stdin:    manager.GetInstance().GetStream(),
	}
	run, err := prompt.Run()
	if err != nil {
		log.Println(err)
	}
	return run

}
