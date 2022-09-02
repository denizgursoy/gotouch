package req

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
	"path/filepath"
	"regexp"
)

type (
	ProjectNameRequirement struct {
		Prompter prompter.Prompter `validate:"required"`
		Manager  manager.Manager   `validate:"required"`
		Logger   logger.Logger     `validate:"required"`
		Store    store.Store       `validate:"required"`
	}

	projectNameTask struct {
		ModuleName string          `validate:"required"`
		Manager    manager.Manager `validate:"required"`
		Logger     logger.Logger   `validate:"required"`
		Store      store.Store     `validate:"required"`
	}
)

const (
	ProjectNameDirection = "Enter Project Name"
)

func (p *ProjectNameRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, nil, err
	}

	moduleName, err := p.Prompter.AskForString(ProjectNameDirection, validateProjectName)

	if err != nil {
		return nil, nil, err
	}

	task := projectNameTask{
		ModuleName: moduleName,
		Manager:    p.Manager,
		Logger:     p.Logger,
		Store:      p.Store,
	}

	tasks := make([]model.Task, 0)
	tasks = append(tasks, &task)

	return tasks, nil, nil
}

func (p *projectNameTask) Complete() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}

	projectName := filepath.Base(p.ModuleName)

	p.Store.SetValue(store.ModuleName, p.ModuleName)
	p.Store.SetValue(store.ProjectName, projectName)

	workingDirectory := p.Manager.GetExtractLocation()
	projectFullPath := fmt.Sprintf("%s/%s", workingDirectory, projectName)
	dirCreationErr := p.Manager.CreateDirectoryIfNotExists(projectFullPath)

	p.Store.SetValue(store.WorkingDirectory, workingDirectory)
	p.Store.SetValue(store.ProjectFullPath, projectFullPath)

	if dirCreationErr != nil {
		return dirCreationErr
	}

	p.Logger.LogInfo(fmt.Sprintf("%s is created", projectFullPath))
	return nil
}

func validateProjectName(name interface{}) error {
	projectName := name.(string)
	compile, err := regexp.Compile("^([a-zA-Z]+(((\\w|(\\.[a-z]+))*)\\/)+[a-zA-Z]+(\\w)*)$|^([a-zA-Z]+\\w*)$")
	if err != nil {
		return errors.New("regex error")
	}
	if compile.MatchString(projectName) {
		return nil
	}
	return errors.New("invalid project name")
}
