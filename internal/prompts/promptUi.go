package prompts

import (
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/manifoldco/promptui"
	"log"
)

type promptUi struct {
	m manager.Manager
}

func (p *promptUi) AskForSelectionFromList(direction string, list []Option) (interface{}, error) {

	count := len(list)

	if count == 0 {
		return nil, ErrProductStructureListIsEmpty
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
		Stdin: p.m.GetStream(),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return list[index], nil
}

func (p *promptUi) AskForString(direction string, validator StringValidator) string {
	prompt := promptui.Prompt{
		Label:    direction,
		Validate: promptui.ValidateFunc(validator),
		Stdin:    p.m.GetStream(),
	}
	run, err := prompt.Run()
	if err != nil {
		log.Println(err)
	}
	return run

}
