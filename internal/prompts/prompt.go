//go:generate mockgen -source=./prompt.go -destination=mock-prompt.go -package=prompts

package prompts

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/manager"
	"sync"
)

var (
	once     = sync.Once{}
	prompter Prompter
)

var (
	ErrProductStructureListIsEmpty = errors.New("options can not be empty")
)

type (
	Prompter interface {
		AskForString(direction string, validator StringValidator) string
		AskForSelectionFromList(direction string, list []Option) (interface{}, error)
	}

	ListOption struct {
		DisplayText string
		ReturnVal   interface{}
	}

	StringValidator func(string) error

	Option interface {
		String() string
	}
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = &promptUi{
			m: manager.GetInstance(),
		}
	})
	return prompter
}
