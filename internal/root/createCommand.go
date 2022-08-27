package root

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/req"
	"github.com/go-playground/validator/v10"
)

var (
	ErrNotYamlFile           = errors.New("file should be a yaml file")
	ErrNotValidUrlOrFilePath = errors.New("file or url is not valid")
	ErrAllFieldsAreRequired  = errors.New("all filed are required")
)

type (
	CreateCommandOptions struct {
		Lister     lister.Lister         `validate:"required"`
		Prompter   prompter.Prompter     `validate:"required"`
		Manager    manager.Manager       `validate:"required"`
		Compressor compressor.Compressor `validate:"required"`
		Executor   executor.Executor     `validate:"required"`
		Logger     logger.Logger         `validate:"required"`
		Path       *string               `validate:"omitempty,endswith=.yaml,url|file"`
	}
)

func CreateNewProject(opts *CreateCommandOptions) error {
	if validationError := isValid(opts); validationError != nil {
		return validationError
	}

	requirements := make(executor.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		Prompter: opts.Prompter,
		Manager:  opts.Manager,
		Logger:   opts.Logger,
	})

	projects, err := opts.Lister.GetProjectList(opts.Path)

	if err != nil {
		return err
	}

	requirements = append(requirements, &req.ProjectStructureRequirement{
		ProjectsData: projects,
		Prompter:     opts.Prompter,
		Compressor:   opts.Compressor,
		Manager:      opts.Manager,
		Logger:       opts.Logger,
		Executor:     opts.Executor,
	})

	return opts.Executor.Execute(requirements)
}

func isValid(opts *CreateCommandOptions) error {
	err := validator.New().Struct(opts)
	if err != nil {
		fieldErrors := err.(validator.ValidationErrors)
		fieldError := fieldErrors[0]
		if fieldError.Field() == "Path" {
			if fieldError.ActualTag() == "endswith" {
				return ErrNotYamlFile
			}
			return ErrNotValidUrlOrFilePath
		}

		return ErrAllFieldsAreRequired
	}
	return nil
}
