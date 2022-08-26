package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/req"
	"github.com/go-playground/validator/v10"
)

type (
	CreateCommandOptions struct {
		Lister     lister.Lister         `validate:"required"`
		Prompter   prompter.Prompter     `validate:"required"`
		Manager    manager.Manager       `validate:"required"`
		Compressor compressor.Compressor `validate:"required"`
		Executor   executor.Executor     `validate:"required"`
		Path       *string               `validate:"omitempty,endswith=.yaml,url|file"`
	}
)

func CreateNewProject(opts *CreateCommandOptions) error {
	if validationError := isValid(opts); validationError != nil {
		return model.ErrMissingField
	}

	requirements := make(executor.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		Prompter: opts.Prompter,
		Manager:  opts.Manager,
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
	})

	return opts.Executor.Execute(requirements)
}

func isValid(opts *CreateCommandOptions) error {
	return validator.New().Struct(opts)
}
