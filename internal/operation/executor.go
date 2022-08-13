package operation

import "errors"

type (
	executor struct {
	}
)

func newExecutor() Executor {
	return executor{}
}

func (e executor) Execute(requirements Requirements) error {

	if requirements == nil {
		return errors.New("req cannot be empty")
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
