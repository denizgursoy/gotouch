package queue

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/model"
)

type Tasks []model.Task

type Requirements []model.Requirement

func Execute(requirements Requirements) error {

	if requirements == nil {
		return errors.New("req cannot be empty")
	}

	tasks := make(Tasks, 0)

	for _, requirement := range requirements {
		task := requirement.AskForInput()
		tasks = append(tasks, task)
	}

	var previousResponse interface{}

	for _, task := range tasks {
		previousResponse = task.Complete(previousResponse)
	}

	return nil
}
