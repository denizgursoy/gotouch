package root

import (
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/req"
)

type (
	CreateNewProjectOptions struct {
	}
)

func CreateNewProject(options *CreateNewProjectOptions) error {
	requirements := make(operation.Requirements, 0)

	requirements = append(requirements, req.ProjectNameRequirement{})

	projects, err := operation.Lister.GetDefaultProjects()

	if err != nil {
		return err
	}

	requirements = append(requirements, req.ProjectStructureRequirement{
		ProjectsData: projects,
	})

	return operation.MainExecutor.Execute(requirements)
}
