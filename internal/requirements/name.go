package requirements

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
	"golang.org/x/mod/module"
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
	ModuleNameDirection = "Enter Module Name   :"
)

func (p *ProjectNameRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, nil, err
	}

	//TODO: move validateModuleName to langs
	moduleName, err := p.Prompter.AskForString(ModuleNameDirection, p.validateModuleName)
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
	dirCreationErr := p.Manager.CreateDirectoryIfNotExist(projectFullPath)

	p.Store.SetValue(store.WorkingDirectory, workingDirectory)
	p.Store.SetValue(store.ProjectFullPath, projectFullPath)

	if dirCreationErr != nil {
		return dirCreationErr
	}

	p.Logger.LogInfo(fmt.Sprintf("%s is created", projectFullPath))
	return nil
}

func (p *ProjectNameRequirement) validateModuleName(name any) error {
	moduleName := name.(string)
	if !strings.Contains(moduleName, "/") {
		err := module.CheckImportPath(moduleName)
		if err != nil {
			return err
		}
	} else {
		err := module.CheckPath(moduleName)
		if err != nil {
			return err
		}
	}

	projectName := filepath.Base(moduleName)
	workingDirectory := p.Manager.GetExtractLocation()
	projectFullPath := fmt.Sprintf("%s/%s", workingDirectory, projectName)
	if _, err := os.Stat(projectFullPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s exists, select another name", projectName)
	}

	return nil
}
