package req

import (
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/go-playground/validator/v10"
)

type (
	dependencyTask struct {
		Manager    manager.Manager `validate:"required"`
		Dependency string          `validate:"required"`
	}
)

func (d *dependencyTask) Complete(i interface{}) (interface{}, error) {
	if err := validator.New().Struct(d); err != nil {
		return nil, err
	}

	if err := d.Manager.AddDependency(d.Dependency); err != nil {
		return nil, err
	}

	return nil, nil
}
