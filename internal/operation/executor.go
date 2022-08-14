package operation

import (
	"errors"
	"log"
	"os/exec"
)

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
		data, err := task.Complete(previousResponse)

		if err != nil {
			return err
		}
		previousResponse = data
	}

	return nil
}

func (e executor) RunCommand(data *CommandData) error {
	cmd := exec.Command(data.Command, data.Args...)
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}
