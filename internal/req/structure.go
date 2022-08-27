package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"path/filepath"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*model.ProjectStructureData `validate:"required"`
		Prompter     prompter.Prompter             `validate:"required"`
		Compressor   compressor.Compressor         `validate:"required"`
		Manager      manager.Manager               `validate:"required"`
		Logger       logger.Logger                 `validate:"required"`
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData `validate:"required"`
		Compressor       compressor.Compressor       `validate:"required"`
		Manager          manager.Manager             `validate:"required"`
		Executor         executor.Executor           `validate:"required"`
		Logger           logger.Logger               `validate:"required"`
	}
)

const (
	SelectProjectTypeDirection = "Select Project Type"
)

func (p *ProjectStructureRequirement) AskForInput() (model.Task, error) {
	options := make([]prompter.Option, 0)
	for _, datum := range p.ProjectsData {
		options = append(options, datum)
	}

	selected, err := p.Prompter.AskForSelectionFromList(SelectProjectTypeDirection, options)

	if err != nil {
		return nil, err
	}

	return &projectStructureTask{
		ProjectStructure: selected.(*model.ProjectStructureData),
		Compressor:       p.Compressor,
		Manager:          p.Manager,
		Logger:           p.Logger,
	}, nil
}

func (p *projectStructureTask) Complete(previousResponse interface{}) (interface{}, error) {
	projectName := previousResponse.(string)
	folderName := filepath.Base(projectName)

	p.Logger.LogInfo("extracting files...")
	if err := p.Compressor.UncompressFromUrl(p.ProjectStructure.URL, folderName); err != nil {
		return nil, err
	}
	p.Logger.LogInfo("zip is extracted successfully")

	p.Logger.LogInfo("changing module name")
	err := p.Manager.EditGoModule(projectName, folderName)
	if err != nil {
		return nil, err
	}
	p.Logger.LogInfo(fmt.Sprintf("module name was changed to %s", projectName))

	return nil, nil
}
