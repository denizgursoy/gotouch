package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*model.ProjectStructureData
		Prompter     prompter.Prompter     `validate:"required"`
		Compressor   compressor.Compressor `validate:"required"`
		Manager      manager.Manager       `validate:"required"`
		Logger       logger.Logger         `validate:"required"`
		Executor     executor.Executor     `validate:"required"`
		Store        store.Store           `validate:"required"`
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData `validate:"required"`
		Compressor       compressor.Compressor       `validate:"required"`
		Manager          manager.Manager             `validate:"required"`
		Executor         executor.Executor           `validate:"required"`
		Logger           logger.Logger               `validate:"required"`
		Store            store.Store                 `validate:"required"`
	}
)

const (
	SelectProjectTypeDirection = "Select Project Type"
)

func (p *ProjectStructureRequirement) AskForInput() (model.Task, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, err
	}

	options := make([]prompter.Option, 0)
	for _, datum := range p.ProjectsData {
		options = append(options, datum)
	}

	selected, err := p.Prompter.AskForSelectionFromList(SelectProjectTypeDirection, options)

	if err != nil {
		return nil, err
	}

	task := projectStructureTask{
		ProjectStructure: selected.(*model.ProjectStructureData),
		Compressor:       p.Compressor,
		Manager:          p.Manager,
		Logger:           p.Logger,
		Executor:         p.Executor,
		Store:            p.Store,
	}

	return &task, nil
}

func (p *projectStructureTask) Complete() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}

	moduleName := p.Store.GetValue(store.ModuleName)
	folderName := p.Store.GetValue(store.ProjectName)

	p.Logger.LogInfo("Extracting files...")
	if err := p.Compressor.UncompressFromUrl(p.ProjectStructure.URL, folderName); err != nil {
		return err
	}
	p.Logger.LogInfo("Zip is extracted successfully")

	p.Logger.LogInfo(fmt.Sprintf("module name will be -> %s", moduleName))
	err := p.Manager.EditGoModule()
	if err != nil {
		return err
	}
	p.Logger.LogInfo(fmt.Sprintf("module name was changed to -> %s", moduleName))

	return nil
}
