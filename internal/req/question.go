package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/langs"
	"net/http"

	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
)

type (
	QuestionRequirement struct {
		Question        model.Question    `validate:"required"`
		Prompter        prompter.Prompter `validate:"required"`
		Logger          logger.Logger     `validate:"required"`
		Executor        executor.Executor `validate:"required"`
		Manager         manager.Manager   `validate:"required"`
		Store           store.Store       `validate:"required"`
		LanguageChecker langs.Checker     `validate:"required"`
	}

	NoneOfAboveChoice struct{}
)

var noneOfAboveChoice = NoneOfAboveChoice{}

func (q *QuestionRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if err := validator.New().Struct(q); err != nil {
		return nil, nil, err
	}

	question := q.Question

	choices := make([]fmt.Stringer, 0)
	for _, choice := range question.Choices {
		choices = append(choices, choice)
	}
	var selection *model.Choice

	if question.CanSkip && len(question.Choices) > 1 {
		choices = append(choices, noneOfAboveChoice)
	}
	isYesNoQuestion := question.CanSkip && len(choices) == 1

	if isYesNoQuestion {
		userSelection, err := q.Prompter.AskForYesOrNo(question.Direction)
		if err != nil {
			return nil, nil, err
		}
		if userSelection {
			selection = question.Choices[0]
		}
	} else {
		userSelection, err := q.Prompter.AskForSelectionFromList(question.Direction, choices)
		if err != nil {
			return nil, nil, err
		}
		if userSelection != noneOfAboveChoice {
			selection = userSelection.(*model.Choice)
		}
	}

	tasks := make([]model.Task, 0)

	if selection != nil {
		for _, dependency := range selection.Dependencies {
			tasks = append(tasks, &dependencyTask{
				Dependency:      dependency,
				LanguageChecker: q.LanguageChecker,
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
		q.Store.StoreValues(selection.Values)
	}
	return tasks, nil, nil
}

func (n NoneOfAboveChoice) String() string {
	return "None of above"
}
