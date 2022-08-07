package prompts

import (
	"sync"
)

var (
	once     = sync.Once{}
	prompter Prompter
)

type (
	Prompter interface {
		AskForString(direction string) string
		AskForSelectionFromList(direction string, listOptions []*ListOption) interface{}
	}

	ListOption struct {
		DisplayText string
		ReturnVal   interface{}
	}
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = promptUi{}
	})
	return prompter
}
