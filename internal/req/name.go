package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/util"
	"log"
	"os"
	"regexp"
	"strings"
)

type ProjectNameRequirement struct {
}

func (p ProjectNameRequirement) AskForInput() (model.Task, error) {
	instance := prompts.GetInstance()
	forString := instance.AskForString("Enter Project Name", validateProjectName)

	if len(strings.TrimSpace(forString)) == 0 {
		return nil, errors.New("project name cannot be empty")
	}

	return projectNameTask{
		ProjectName: forString,
	}, nil
}

type projectNameTask struct {
	ProjectName string
}

func (p projectNameTask) Complete(interface{}) interface{} {
	path, _ := util.GetBaseName(p.ProjectName)
	directoryPath := fmt.Sprintf("./%s", path)
	_ = os.Mkdir(directoryPath, os.ModePerm)
	return p.ProjectName
}

func validateProjectName(projectName string) error {
	compile, err := regexp.Compile("^[a-zA-Z](\\w|\\.|_|-)*((\\/)?[a-zA-Z](\\w|\\.|_|-)*)*$")
	if err != nil {
		log.Fatalln("regex error")
	}
	if compile.MatchString(projectName) {
		return nil
	}
	return errors.New("invalid project name")
}
