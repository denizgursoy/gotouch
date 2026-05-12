package executor

import (
	"context"
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

func (e executor) Execute(ctx context.Context, requirements Requirements) error {
	return e.ExecuteWithHook(ctx, requirements, nil)
}

func (e executor) ExecuteWithHook(ctx context.Context, requirements Requirements, onPromptsDone func()) error {
	if requirements == nil {
		return EmptyRequirementError
	}

	tasks := make(Tasks, 0)
	var promptErr error

	for i := 0; i < len(requirements); i++ {
		inputTask, inputRequirements, err := requirements[i].AskForInput()
		if err != nil {
			promptErr = err
			break
		}
		tasks = append(tasks, inputTask...)
		requirements = append(requirements, inputRequirements...)
	}

	if onPromptsDone != nil {
		onPromptsDone()
	}

	if promptErr != nil {
		return promptErr
	}

	for _, task := range tasks {
		err := task.Complete(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
