package root

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/req"
)

type (
	CreateCommandOptions struct {
		lister     lister.Lister
		prompter   prompter.Prompter
		manager    manager.Manager
		compressor compressor.Compressor
		executor   executor.Executor
	}
)

var (
	ErrMissingField = errors.New("all fields should be provided")
)

func CreateNewProject(opts *CreateCommandOptions) error {

	if !isValid(opts) {
		return ErrMissingField
	}

	requirements := make(executor.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		Prompter: opts.prompter,
		Manager:  opts.manager,
	})

	projects, err := opts.lister.GetProjectList(nil)

	if err != nil {
		return err
	}

	requirements = append(requirements, &req.ProjectStructureRequirement{
		ProjectsData: projects,
		Prompter:     opts.prompter,
		Compressor:   opts.compressor,
		Manager:      opts.manager,
	})

	return opts.executor.Execute(requirements)
}

func isValid(opts *CreateCommandOptions) bool {
	return opts.compressor != nil &&
		opts.executor != nil &&
		opts.lister != nil &&
		opts.prompter != nil &&
		opts.manager != nil
}
