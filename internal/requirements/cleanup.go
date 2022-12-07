package requirements

import (
	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/go-playground/validator/v10"
)

type (
	cleanupRequirement struct {
		LanguageChecker langs.Checker `validate:"required"`
	}
	cleanupTask struct {
		LanguageChecker langs.Checker `validate:"required"`
	}
)

func (c *cleanupRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	tasks := make([]model.Task, 0)

	tasks = append(tasks, &cleanupTask{
		LanguageChecker: c.LanguageChecker,
	})
	return tasks, nil, nil
}

func (c *cleanupTask) Complete() error {
	if err := validator.New().Struct(c); err != nil {
		return err
	}
	return c.LanguageChecker.CleanUp()
}
