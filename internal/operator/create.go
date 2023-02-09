package operator

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/config"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/requirements"
	"github.com/denizgursoy/gotouch/internal/store"
)

var (
	ErrNotYamlFile           = errors.New("file should be a yaml file")
	ErrNotValidUrlOrFilePath = errors.New("file or url is not valid")
	ErrAllFieldsAreRequired  = errors.New("all filed are required")
)

type (
	CreateNewProjectOptions struct {
		Lister        lister.Lister         `validate:"required"`
		Prompter      prompter.Prompter     `validate:"required"`
		Manager       manager.Manager       `validate:"required"`
		Compressor    compressor.Compressor `validate:"required"`
		Executor      executor.Executor     `validate:"required"`
		Logger        logger.Logger         `validate:"required"`
		Path          *string               `validate:"omitempty,endswith=.yaml,url|file"`
		Store         store.Store           `validate:"required"`
		Cloner        cloner.Cloner         `validate:"required"`
		CommandRunner commandrunner.Runner  `validate:"required"`
		ConfigManager config.ConfigManager  `validate:"required"`
	}
)

func (o *operator) CreateNewProject(opts *CreateNewProjectOptions) error {
	if validationError := isValid(opts); validationError != nil {
		return validationError
	}

	newProjectRequirements := make(executor.Requirements, 0)
	if opts.Path == nil || len(strings.TrimSpace(*opts.Path)) == 0 {
		path, err := opts.ConfigManager.GetDefaultPath()
		if err != nil {
			return err
		}
		opts.Path = &path
	}

	projects, err := opts.Lister.GetProjectList(opts.Path)
	if err != nil {
		return err
	}

	requirement := requirements.ProjectStructureRequirement{
		ProjectsData:  projects,
		Prompter:      opts.Prompter,
		Compressor:    opts.Compressor,
		Manager:       opts.Manager,
		Logger:        opts.Logger,
		Executor:      opts.Executor,
		Store:         opts.Store,
		Cloner:        opts.Cloner,
		CommandRunner: opts.CommandRunner,
	}
	newProjectRequirements = append(newProjectRequirements, &requirement)

	return opts.Executor.Execute(newProjectRequirements)
}

func isValid(opts *CreateNewProjectOptions) error {
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
