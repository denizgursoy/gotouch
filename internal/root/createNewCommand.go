package root

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/req"
)

type (
	CreateNewProjectOptions struct {
		l lister.Lister
		p prompts.Prompter
		m manager.Manager
	}
)

func CreateNewProject(options *CreateNewProjectOptions) error {
	requirements := make(operation.Requirements, 0)

	requirements = append(requirements, &req.ProjectNameRequirement{
		P: options.p,
		M: options.m,
	})

	projects, err := options.l.GetDefaultProjects()

	if err != nil {
		return err
	}

	requirements = append(requirements, req.ProjectStructureRequirement{
		ProjectsData: projects,
	})

	return operation.GetInstance().Execute(requirements)
}
