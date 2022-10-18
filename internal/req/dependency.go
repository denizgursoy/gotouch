package req

import (
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
)

type (
	dependencyTask struct {
		Dependency      interface{}       `validate:"required"`
		Logger          logger.Logger     `validate:"required"`
		Executor        executor.Executor `validate:"required"`
		Store           store.Store       `validate:"required"`
		LanguageChecker langs.Checker     `validate:"required"`
	}
)

func (d *dependencyTask) Complete() error {
	if err := validator.New().Struct(d); err != nil {
		return err
	}
	return d.LanguageChecker.GetLangChecker().GetDependency(d.Dependency)
}
