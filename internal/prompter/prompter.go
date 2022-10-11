//go:generate mockgen -source=./prompter.go -destination=mockPrompter.go -package=prompter

package prompter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/denizgursoy/gotouch/internal/manager"
)

var (
	once      = sync.Once{}
	prompter  Prompter
	EmptyList = errors.New("options can not be empty")
)

type (
	Prompter interface {
		AskForString(direction string, validator Validator) (string, error)
		AskForSelectionFromList(direction string, list []fmt.Stringer) (interface{}, error)
		AskForYesOrNo(direction string) (bool, error)
		AskForMultilineString(direction, defaultValue, pattern string) (string, error)
	}

	ListOption struct {
		DisplayText string
		ReturnVal   interface{}
	}

	Validator func(interface{}) error

	srv struct {
		Manager manager.Manager
	}
)

func GetInstance() Prompter {
	once.Do(func() {
		prompter = &srv{
			Manager: manager.GetInstance(),
		}
	})
	return prompter
}
