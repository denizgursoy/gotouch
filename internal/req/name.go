package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/go-playground/validator/v10"
	"path/filepath"
	"regexp"
)

type (
	ProjectNameRequirement struct {
		Prompter prompter.Prompter `validate:"required"`
		Manager  manager.Manager   `validate:"required"`
		Logger   logger.Logger     `validate:"required"`
	}

	projectNameTask struct {
		ProjectName string          `validate:"required"`
		Manager     manager.Manager `validate:"required"`
		Logger      logger.Logger   `validate:"required"`
	}
)

const (
	ProjectNameDirection = "Enter Project Name"
)

func (p *ProjectNameRequirement) AskForInput() (model.Task, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, err
	}

	projectName, err := p.Prompter.AskForString(ProjectNameDirection, validateProjectName)

	if err != nil {
		return nil, err
	}

	return &projectNameTask{
		ProjectName: projectName,
		Manager:     p.Manager,
		Logger:      p.Logger,
	}, nil
}

func (p *projectNameTask) Complete(interface{}) (interface{}, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, err
	}

	folderName := filepath.Base(p.ProjectName)
	directoryPath := fmt.Sprintf("%s/%s", p.Manager.GetExtractLocation(), folderName)
	dirCreationErr := p.Manager.CreateDirectoryIfNotExists(directoryPath)

	if dirCreationErr != nil {
		return nil, dirCreationErr
	}

	p.Logger.LogInfo(fmt.Sprintf("%s is created", directoryPath))
	return p.ProjectName, nil
}

func validateProjectName(projectName string) error {
	compile, err := regexp.Compile("^([a-zA-Z]+(((\\w|(\\.[a-z]+))*)\\/)+[a-zA-Z]+(\\w)*)$|^([a-zA-Z]+\\w*)$")
	if err != nil {
		return errors.New("regex error")
	}
	if compile.MatchString(projectName) {
		return nil
	}
	return errors.New("invalid project name")
}
