package requirements

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/mod/module"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
)

type (
	ProjectNameRequirement struct {
		Prompter     prompter.Prompter `validate:"required"`
		Manager      manager.Manager   `validate:"required"`
		Logger       logger.Logger     `validate:"required"`
		Store        store.Store       `validate:"required"`
		InitialValue string
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

	// TODO: move validateModuleName to langs
	moduleName, err := p.Prompter.AskForString(ModuleNameDirection, p.InitialValue, p.validateModuleName)
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

func (p *projectNameTask) Complete(ctx context.Context) error {
	if err := validator.New().StructCtx(ctx, p); err != nil {
		return err
	}

	projectName := filepath.Base(strings.TrimSpace(p.ModuleName))
	projectName = sanitizeProjectName(projectName)
	p.Store.SetValue(store.ModuleName, p.ModuleName)
	p.Store.SetValue(store.ProjectName, projectName)

	inline, inlineParseError := strconv.ParseBool(p.Store.GetValue(store.Inline))
	if inlineParseError != nil {
		return inlineParseError
	}

	workingDirectory := p.Manager.GetExtractLocation()
	projectFullPath := workingDirectory
	if !inline {
		projectFullPath = filepath.Join(workingDirectory, projectName)

		dirCreationErr := p.Manager.CreateDirectoryIfNotExist(projectFullPath)
		if dirCreationErr != nil {
			return dirCreationErr
		}
	}

	p.Store.SetValue(store.WorkingDirectory, workingDirectory)
	p.Store.SetValue(store.ProjectFullPath, projectFullPath)

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

	inline, err := strconv.ParseBool(p.Store.GetValue(store.Inline))
	if err != nil {
		return err
	}

	if inline {
		return nil
	}

	projectName := sanitizeProjectName(filepath.Base(moduleName))
	workingDirectory := p.Manager.GetExtractLocation()
	projectFullPath := filepath.Join(workingDirectory, projectName)
	if _, err = os.Stat(projectFullPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory %s exists, select another name", projectName)
	}

	return nil
}

func sanitizeProjectName(projectName string) string {
	suffix := strings.TrimSuffix(projectName, ".git")
	suffix = strings.ReplaceAll(suffix, " ", "_")

	return suffix
}
