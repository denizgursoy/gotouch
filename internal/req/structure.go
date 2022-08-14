package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"log"
	"os"
	"path/filepath"
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

func (p projectStructureTask) Complete(previousResponse interface{}) (interface{}, error) {
	projectName := previousResponse.(string)
	folderName := filepath.Base(projectName)

	operation.Extractor.UncompressFromUrl(p.ProjectStructure.URL, folderName)
	return nil, editGoModule(projectName, folderName)
}

func editGoModule(projectName, folderName string) error {
	workingDirectory, err := os.Getwd()

	if err != nil {
		log.Println(err)
		return err
	}

	projectDirectory := fmt.Sprintf("%s/%s", workingDirectory, folderName)

	err = os.Chdir(projectDirectory)
	if err != nil {
		log.Println(err)
		return err
	}

	args := make([]string, 0)

	if hasGoModule(projectDirectory) {
		args = append(args, "mod", "edit", "-module", projectName)
	} else {
		args = append(args, "mod", "init", projectName)
	}

	data := &operation.CommandData{
		Command: "go",
		Args:    args,
	}

	return operation.MainExecutor.RunCommand(data)
}

func hasGoModule(projectDirectory string) bool {
	path := fmt.Sprintf("%s/go.mod", projectDirectory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
