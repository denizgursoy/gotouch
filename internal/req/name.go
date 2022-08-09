package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/util"
	"os"
	"strings"
)

type ProjectNameRequirement struct {
}

func (p ProjectNameRequirement) AskForInput() (model.Task, error) {
	instance := prompts.GetInstance()
	forString := instance.AskForString("Enter Project Name")

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
