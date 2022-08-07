package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"os"
)

type ProjectNameRequirement struct {
}

func (p ProjectNameRequirement) AskForInput() model.Task {
	instance := prompts.GetInstance()
	forString := instance.AskForString("Enter Project Name")
	return projectNameTask{
		ProjectName: forString,
	}
}

type projectNameTask struct {
	ProjectName string
}

func (p projectNameTask) Complete(interface{}) interface{} {
	directoryPath := fmt.Sprintf("./%s", p.ProjectName)
	_ = os.Mkdir(directoryPath, os.ModePerm)
	return p.ProjectName
}
