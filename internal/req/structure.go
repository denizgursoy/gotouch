package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
	"log"
	"os"
	"path/filepath"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*model.ProjectStructureData
		P            prompts.Prompter
		U            uncompressor.Uncompressor
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData
		U                uncompressor.Uncompressor
	}
)

const (
	SelectProjectTypeDirection = "Select Project Type"
)

func (p *ProjectStructureRequirement) AskForInput() (model.Task, error) {

	options := make([]prompts.Option, 0)
	for _, datum := range p.ProjectsData {
		options = append(options, datum)
	}

	selected, err := p.P.AskForSelectionFromList(SelectProjectTypeDirection, options)

	if err != nil {
		return nil, err
	}

	return &projectStructureTask{
		ProjectStructure: selected.(*model.ProjectStructureData),
		U:                p.U,
	}, nil
}

func (p *projectStructureTask) Complete(previousResponse interface{}) (interface{}, error) {
	projectName := previousResponse.(string)
	folderName := filepath.Base(projectName)

	p.U.UncompressFromUrl(p.ProjectStructure.URL, folderName)
	return nil, editGoModule(projectName, folderName)
}

func editGoModule(projectName, folderName string) error {
	workingDirectory := manager.GetInstance().GetExtractLocation()
	projectDirectory := fmt.Sprintf("%s/%s", workingDirectory, folderName)
	fmt.Println(projectDirectory, "projectDirectory")
	fmt.Println(hasGoModule(projectDirectory), "hasGoModule(projectDirectory)")
	fmt.Println(manager.GetInstance().IsTest(), "IsTest")

	err := os.Chdir(projectDirectory)
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

	return operation.GetInstance().RunCommand(data)
}

func hasGoModule(projectDirectory string) bool {
	path := fmt.Sprintf("%s/go.mod", projectDirectory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
