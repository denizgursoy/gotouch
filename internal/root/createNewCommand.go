package root

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/req"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
)

type (
	CreateNewProjectOptions struct {
		lister       lister.Lister
		prompter     prompts.Prompter
		manager      manager.Manager
		uncompressor uncompressor.Uncompressor
		executor     operation.Executor
	}
)

func CreateNewProject(opts *CreateNewProjectOptions) error {
	requirements := make(operation.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		Prompter: opts.prompter,
		Manager:  opts.manager,
	})

	projects, err := opts.lister.GetDefaultProjects()

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
