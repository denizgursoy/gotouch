package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/go-playground/validator/v10"
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

	NoneOfAboveOption struct{}
)

var (
	noneOfAboveOption = NoneOfAboveOption{}
)

func (q *QuestionRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if err := validator.New().Struct(q); err != nil {
		return nil, nil, err
	}

	question := q.Question

	options := make([]fmt.Stringer, 0)
	for _, option := range question.Options {
		options = append(options, option)
	}
	var selection *model.Option

	if question.CanSkip && len(question.Options) > 1 {
		options = append(options, noneOfAboveOption)
	}
	isYesNoQuestion := question.CanSkip && len(options) == 1

	if isYesNoQuestion {
		userSelection, err := q.Prompter.AskForYesOrNo(question.Direction)
		if err != nil {
			return nil, nil, err
		}
		if userSelection {
			selection = question.Options[0]
		}
	} else {
		userSelection, err := q.Prompter.AskForSelectionFromList(question.Direction, options)
		if err != nil {
			return nil, nil, err
		}
		if userSelection != noneOfAboveOption {
			selection = userSelection.(*model.Option)
		}
	}

	tasks := make([]model.Task, 0)

	if selection != nil {
		for _, dependency := range selection.Dependencies {
			tasks = append(tasks, &dependencyTask{
				Dependency: *dependency,
				Logger:     q.Logger,
				Executor:   q.Executor,
			})
		}

		for _, file := range selection.Files {
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

func (n NoneOfAboveOption) String() string {
	return "None of above"
}
