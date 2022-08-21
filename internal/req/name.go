package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"log"
	"path/filepath"
	"regexp"
)

type (
	ProjectNameRequirement struct {
		P prompts.Prompter
		M manager.Manager
	}

	projectNameTask struct {
		ProjectName string
		m           manager.Manager
	}
)

func (p *ProjectNameRequirement) AskForInput() (model.Task, error) {
	projectName := p.P.AskForString("Enter Project Name", validateProjectName)

	return &projectNameTask{
		ProjectName: projectName,
		m:           p.M,
	}, nil
}

func (p *projectNameTask) Complete(interface{}) (interface{}, error) {
	folderName := filepath.Base(p.ProjectName)
	directoryPath := fmt.Sprintf("%s/%s", p.m.GetExtractLocation(), folderName)
	dirCreationErr := p.m.CreateDirectoryIfNotExists(directoryPath)
	return p.ProjectName, dirCreationErr
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
