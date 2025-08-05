package requirements

import (
	"github.com/denizgursoy/gotouch/internal/langs"

	"context"

	"github.com/go-playground/validator/v10"
)

type (
	dependencyTask struct {
		Dependency      any           `validate:"required"`
		LanguageChecker langs.Checker `validate:"required"`
	}
)

func (d *dependencyTask) Complete(ctx context.Context) error {
	if err := validator.New().StructCtx(ctx, d); err != nil {
		return err
	}
	return d.LanguageChecker.GetDependency(d.Dependency)
}
