package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/req"
)

type (
	CreateCommandOptions struct {
		lister       lister.Lister
		prompter     prompter.Prompter
		manager      manager.Manager
		uncompressor compressor.Uncompressor
		executor     executor.Executor
	}
)

func CreateNewProject(opts *CreateCommandOptions) error {
	requirements := make(executor.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		Prompter: opts.prompter,
		Manager:  opts.manager,
	})

	projects, err := opts.lister.GetProjectList()

	if err != nil {
		return err
	}

	requirements = append(requirements, &req.ProjectStructureRequirement{
		ProjectsData: projects,
		Prompter:     opts.prompter,
		Uncompressor: opts.uncompressor,
	})

	return opts.executor.Execute(requirements)
}
