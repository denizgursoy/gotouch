package prompts

import (
	"github.com/manifoldco/promptui"
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
	}

	index, _, _ := prompt.Run()

	return listOptions[index].ReturnVal
}

func (p promptUi) AskForString(direction string) string {
	prompt := promptui.Prompt{
		Label: direction,
	}
	run, _ := prompt.Run()
	return run

}
