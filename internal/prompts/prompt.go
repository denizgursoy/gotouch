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
		AskForString(direction string, validator StringValidator) string
		AskForSelectionFromList(direction string, listOptions []*ListOption) interface{}
	}

	ListOption struct {
		DisplayText string
		ReturnVal   interface{}
	}

	StringValidator func(string) error
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = promptUi{}
	})
	return prompter
}
