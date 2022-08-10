package operation

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
)

type (
	Tasks []model.Task

	Requirements []model.Requirement
)

var (
	Prompter  = prompts.GetInstance()
	Extractor = uncompressor.GetInstance()
	Lister    = lister.GetInstance()
)

func Execute(requirements Requirements) error {

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
