package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/operation"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type ProjectNameRequirement struct {
}

func (p ProjectNameRequirement) AskForInput() (model.Task, error) {

	projectName := operation.Prompter.AskForString("Enter Project Name", validateProjectName)

	return projectNameTask{
		ProjectName: projectName,
	}, nil
}

type projectNameTask struct {
	ProjectName string
}

func (p projectNameTask) Complete(interface{}) (interface{}, error) {
	folderName := filepath.Base(p.ProjectName)
	directoryPath := fmt.Sprintf("%s/%s", prompts.GetExtractLocation(), folderName)
	err := os.Mkdir(directoryPath, os.ModePerm)
	return p.ProjectName, err
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
