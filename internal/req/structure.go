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
		Uncompressor compressor.Uncompressor
		Manager      manager.Manager
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData
		Uncompressor     compressor.Uncompressor
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
		Uncompressor:     p.Uncompressor,
		Manager:          p.Manager,
	}, nil
}

func (p *projectStructureTask) Complete(previousResponse interface{}) (interface{}, error) {
	projectName := previousResponse.(string)
	folderName := filepath.Base(projectName)

	p.Uncompressor.UncompressFromUrl(p.ProjectStructure.URL, folderName)
	return nil, p.Manager.EditGoModule(projectName, folderName)
}
