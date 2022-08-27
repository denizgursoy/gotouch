package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/go-playground/validator/v10"
	"strings"
)

type (
	dependencyTask struct {
		Dependency string            `validate:"required"`
		Logger     logger.Logger     `validate:"required"`
		Executor   executor.Executor `validate:"required"`
	}
)

var (
	atSign                  = "@"
	latestVersion           = "latest"
	ErrDependencyIsNotValid = errors.New("dependency is not valid")
)

func (d *dependencyTask) Complete(i interface{}) (interface{}, error) {
	if err := validator.New().Struct(d); err != nil {
		return nil, err
	}

	if len(strings.TrimSpace(d.Dependency)) == 0 {
		return nil, ErrDependencyIsNotValid
	}

	dependency := manager.Dependency{
		Url:     &d.Dependency,
		Version: &latestVersion,
	}

	atIndex := strings.LastIndex(d.Dependency, atSign)
	if atIndex != -1 {
		url := d.Dependency[:atIndex]
		version := d.Dependency[atIndex+1:]

		dependency.Url = &url
		dependency.Version = &version
	}

	d.Logger.LogInfo(fmt.Sprintf("Adding dependency -> %s", dependency.String()))

	data := &executor.CommandData{
		Command: "go",
		Args:    []string{"get", dependency.String()},
	}

	if err := d.Executor.RunCommand(data); err != nil {
		return nil, err
	}

	d.Logger.LogInfo(fmt.Sprintf("Dependency add succesfully -> %s", dependency.String()))

	return nil, nil
}
