package req

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"path/filepath"
)

type (
	ProjectStructureRequirement struct {
		ProjectsData []*model.ProjectStructureData
		Prompter     prompter.Prompter
		Compressor   compressor.Compressor
		Manager      manager.Manager
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData
		Compressor       compressor.Compressor
		Manager          manager.Manager
		Executor         executor.Executor
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
	}, nil
}

func (p *projectStructureTask) Complete(previousResponse interface{}) (interface{}, error) {
	projectName := previousResponse.(string)
	folderName := filepath.Base(projectName)

	if err := p.Compressor.UncompressFromUrl(p.ProjectStructure.URL, folderName); err != nil {
		return nil, err
	}

	return nil, p.Manager.EditGoModule(projectName, folderName)
}
