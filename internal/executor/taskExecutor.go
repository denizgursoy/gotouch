package executor

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/store"
)

type (
	executor struct {
		Store store.Store `validate:"required"`
	}
)

var EmptyRequirementError = errors.New("requirements cannot be empty")

func newExecutor() Executor {
	return &executor{
		Store: store.GetInstance(),
	}
}

func (e executor) Execute(requirements Requirements) error {
	if requirements == nil {
		return EmptyRequirementError
	}

	tasks := make(Tasks, 0)

	for i := 0; i < len(requirements); i++ {
		inputTask, inputRequirements, err := requirements[i].AskForInput()
		if err != nil {
			return err
		}
		tasks = append(tasks, inputTask...)
		requirements = append(requirements, inputRequirements...)
	}

	for _, task := range tasks {
		err := task.Complete()
		if err != nil {
			return err
		}
	}

	return nil
}
