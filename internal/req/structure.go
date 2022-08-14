package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/util"
	"io/ioutil"
	"log"
	"os"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*lister.ProjectStructureData
	}

	projectStructureTask struct {
		ProjectStructure *lister.ProjectStructureData
	}
)

func (p ProjectStructureRequirement) AskForInput() (model.Task, error) {

	instance := prompts.GetInstance()

	projectList := make([]*prompts.ListOption, 0)

	for _, project := range p.ProjectsData {
		displayText := fmt.Sprintf("%s (%s)", project.Name, project.Reference)
		projectList = append(projectList, &prompts.ListOption{
			DisplayText: displayText,
			ReturnVal:   project,
		})
	}

	selected := instance.AskForSelectionFromList("select project type", projectList).(*lister.ProjectStructureData)
	return projectStructureTask{
		ProjectStructure: selected,
	}, nil
}

func (p projectStructureTask) Complete(previousResponse interface{}) interface{} {
	projectName := previousResponse.(string)
	path, err := util.GetBaseName(projectName)
	if err != nil {
		log.Printf("%v", err)
	}

	operation.Extractor.UncompressFromUrl(p.ProjectStructure.URL, path)
	editGoModule(projectName, path)
	return nil
}

func hasGoModule(path string) bool {
	path = fmt.Sprintf("./%s/", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}
	return false
}

func editGoModule(projectName, path string) {

	dir := fmt.Sprintf("./%s", path)
	err := os.Chdir(dir)
	if err != nil {
		log.Printf("Changed Directory error: %v", err)
	}

	args := make([]string, 0)

	if hasGoModule(path) {
		args = append(args, "mod", "edit", "-module", projectName)
	} else {
		args = append(args, "mod", "init", projectName)
	}

	data := &operation.CommandData{
		Command: "go",
		Args:    args,
	}

	err = operation.MainExecutor.RunCommand(data)
}
