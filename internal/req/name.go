package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"log"
	"path/filepath"
	"regexp"
)

type (
	ProjectNameRequirement struct {
		Prompter prompter.Prompter
		Manager  manager.Manager
	}

	projectNameTask struct {
		ProjectName string
		m           manager.Manager
	}
)

func (p *ProjectNameRequirement) AskForInput() (model.Task, error) {
	projectName, err := p.Prompter.AskForString("Enter Project Name", validateProjectName)

	if err != nil {
		return nil, err
	}

	return &projectNameTask{
		ProjectName: projectName,
		m:           p.Manager,
	}, nil
}

func (p *projectNameTask) Complete(interface{}) (interface{}, error) {
	folderName := filepath.Base(p.ProjectName)
	directoryPath := fmt.Sprintf("%s/%s", p.m.GetExtractLocation(), folderName)
	dirCreationErr := p.m.CreateDirectoryIfNotExists(directoryPath)

	if dirCreationErr != nil {
		return nil, dirCreationErr
	}
	return p.ProjectName, nil
}

func validateProjectName(projectName string) error {
	compile, err := regexp.Compile("^([a-zA-Z]+(((\\w|(\\.[a-z]+))*)\\/)+[a-zA-Z]+(\\w)*)$|^([a-zA-Z]+\\w*)$")
	if err != nil {
		log.Fatalln("regex error")
	}
	if compile.MatchString(projectName) {
		return nil
	}
	return errors.New("invalid project name")
}
