package req

import (
	"fmt"
	extractor2 "github.com/denizgursoy/gotouch/internal/extractor"
	lister2 "github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*lister2.ProjectStructureData
	}

	projectStructureTask struct {
		ProjectStructure *lister2.ProjectStructureData
	}
)

func (p ProjectStructureRequirement) AskForInput() model.Task {

	instance := prompts.GetInstance()

	projectList := make([]*prompts.ListOption, 0)

	for _, project := range p.ProjectsData {
		displayText := fmt.Sprintf("%s (%s)", project.Name, project.Reference)
		projectList = append(projectList, &prompts.ListOption{
			DisplayText: displayText,
			ReturnVal:   project,
		})
	}

	selected := instance.AskForSelectionFromList("select project type", projectList).(*lister2.ProjectStructureData)
	return projectStructureTask{
		ProjectStructure: selected,
	}
}

func (p projectStructureTask) Complete(previousResponse interface{}) interface{} {
	projectName := previousResponse.(string)
	ex := extractor2.GetInstance()
	ex.Extract(p.ProjectStructure.URL, projectName)
	return nil
}
