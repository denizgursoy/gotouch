package requirements

import (
	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/go-playground/validator/v10"
)

type (
	cleanupTask struct {
		LanguageChecker langs.Checker `validate:"required"`
	}
)

func (c *cleanupTask) Complete() error {
	if err := validator.New().Struct(c); err != nil {
		return err
	}
	return c.LanguageChecker.CleanUp()
}
