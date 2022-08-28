package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"net/http"
)

type (
	QuestionRequirement struct {
		Question model.Question    `validate:"required"`
		Prompter prompter.Prompter `validate:"required"`
		Logger   logger.Logger     `validate:"required"`
		Executor executor.Executor `validate:"required"`
		Manager  manager.Manager   `validate:"required"`
	}

	NoOption struct{}

	NoneOfAboveOption struct{}
)

var (
	noOption          = NoOption{}
	noneOfAboveOption = NoneOfAboveOption{}
)

func (q *QuestionRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	question := q.Question

	options := make([]fmt.Stringer, 0)
	for _, option := range question.Options {
		options = append(options, option)
	}

	if question.CanSkip {
		if len(question.Options) > 1 {
			options = append(options, noneOfAboveOption)
		} else {
			options = append(options, noOption)
		}
	}

	selection, err := q.Prompter.AskForSelectionFromList(question.Direction, options)
	if err != nil {
		return nil, nil, err
	}

	tasks := make([]model.Task, 0)

	if selection != noOption && selection != noneOfAboveOption {
		selectedOption := selection.(*model.Option)

		for _, dependency := range selectedOption.Dependencies {
			tasks = append(tasks, &dependencyTask{
				Dependency: *dependency,
				Logger:     q.Logger,
				Executor:   q.Executor,
			})
		}

		for _, file := range selectedOption.Files {
			tasks = append(tasks, &fileTask{
				File:    *file,
				Logger:  q.Logger,
				Manager: q.Manager,
				Client:  &http.Client{},
			})
		}

	}
	return tasks, nil, nil
}

func (n NoOption) String() string {
	return "No"
}

func (n NoneOfAboveOption) String() string {
	return "None of above"
}
