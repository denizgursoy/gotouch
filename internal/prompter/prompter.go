//go:generate mockgen -source=./prompter.go -destination=mockPrompter.go -package=prompter --typed

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
		AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error)
		AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error)
		AskForYesOrNo(direction string) (bool, error)
		AskForMultilineString(direction, defaultValue, pattern string) (string, error)
	}

	ListOption struct {
		DisplayText string
		ReturnVal   any
	}

	Validator func(any) error

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
