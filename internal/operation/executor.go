package operation

import "errors"

type (
	executor struct {
	}
)

var (
	EmptyRequirementError = errors.New("requirements cannot be empty")
)

func newExecutor() Executor {
	return executor{}
}

func (e executor) Execute(requirements Requirements) error {

	if requirements == nil {
		return EmptyRequirementError
	}

	tasks := make(Tasks, 0)

	for _, requirement := range requirements {
		task, _ := requirement.AskForInput()
		tasks = append(tasks, task)
	}

	var previousResponse interface{}

	for _, task := range tasks {
		previousResponse = task.Complete(previousResponse)
	}

	return nil
}
