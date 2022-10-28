package requirements

import (
	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/go-playground/validator/v10"
)

type (
	dependencyTask struct {
		Dependency      interface{}   `validate:"required"`
		LanguageChecker langs.Checker `validate:"required"`
	}
)

func (d *dependencyTask) Complete() error {
	if err := validator.New().Struct(d); err != nil {
		return err
	}
	return d.LanguageChecker.GetDependency(d.Dependency)
}
