package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/queue"
	"github.com/denizgursoy/gotouch/internal/util"
	"log"
	"os"
	"regexp"
)

type ProjectNameRequirement struct {
}

func (p ProjectNameRequirement) AskForInput() (model.Task, error) {

	projectName := queue.Prompter.AskForString("Enter Project Name", validateProjectName)

	return projectNameTask{
		ProjectName: projectName,
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
