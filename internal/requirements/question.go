package requirements

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

	if question.CanSkip && len(question.Choices) > 1 {
		choices = append(choices, noneOfAboveChoice)
	}

	isYesNoQuestion := question.CanSkip && len(choices) == 1
	selectedChoices := make([]*model.Choice, 0)

	if isYesNoQuestion {
		userSelection, err := q.Prompter.AskForYesOrNo(question.Direction)
		if err != nil {
			return nil, nil, err
		}
		if userSelection {
			selectedChoices = append(selectedChoices, question.Choices[0])
		}
	} else if question.CanSelectMultiple {
		allSelectedChoices, err := q.Prompter.AskForMultipleSelectionFromList(question.Direction, choices)
		if err != nil {
			return nil, nil, err
		}
		for _, selectedChoice := range allSelectedChoices {
			selectedChoices = append(selectedChoices, selectedChoice.(*model.Choice))
		}
	} else {
		userSelection, err := q.Prompter.AskForSelectionFromList(question.Direction, choices)
		if err != nil {
			return nil, nil, err
		}
		if userSelection != noneOfAboveChoice {
			selectedChoices = append(selectedChoices, userSelection.(*model.Choice))
		}
	}

	tasks := make([]model.Task, 0)

	for _, selection := range selectedChoices {
		choiceTasks := getTasks(selection.Resources, q.Logger, q.Manager, q.LanguageChecker, q.Store)
		tasks = append(tasks, choiceTasks...)
		q.Store.AddValues(selection.Values)
	}

	return tasks, nil, nil
}

func (n NoneOfAboveChoice) String() string {
	return "None of above"
}

func getTasks(r model.Resources, l logger.Logger, m manager.Manager,
	lc langs.Checker, s store.Store) []model.Task {

	tasks := make([]model.Task, 0)

	for _, dependency := range r.Dependencies {
		tasks = append(tasks, &dependencyTask{
			Dependency:      dependency,
			LanguageChecker: lc,
		})
		s.AddDependency(dependency)
	}

	for _, file := range r.Files {
		tasks = append(tasks, &fileTask{
			File:    *file,
			Logger:  l,
			Manager: m,
			Client:  &http.Client{},
		})
	}
	return tasks
}
