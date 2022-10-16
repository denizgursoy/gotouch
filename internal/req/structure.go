package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/langs"

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

func (p *ProjectStructureRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if err := validator.New().Struct(p); err != nil {
		return nil, nil, err
	}

	options := make([]fmt.Stringer, 0)
	for _, datum := range p.ProjectsData {
		options = append(options, datum)
	}

	selected, err := p.Prompter.AskForSelectionFromList(SelectProjectTypeDirection, options)
	if err != nil {
		return nil, nil, err
	}

	p2 := &ProjectNameRequirement{
		Prompter: p.Prompter,
		Manager:  p.Manager,
		Logger:   p.Logger,
		Store:    p.Store,
	}

	nameTasks, nameRequirements, err := p2.AskForInput()
	if err != nil {
		return nil, nil, err
	}

	tasks := make([]model.Task, 0)
	requirements := make([]model.Requirement, 0)

	for _, task := range nameTasks {
		tasks = append(tasks, task)
	}

	for _, requirement := range nameRequirements {
		requirements = append(requirements, requirement)
	}

	projectStructureData := selected.(*model.ProjectStructureData)

	p.setLanguageChecker(projectStructureData)

	if setupError := p.Store.GetLanguageChecker().CheckSetup(); setupError != nil {
		return nil, nil, setupError
	}

	task := projectStructureTask{
		ProjectStructure: projectStructureData,
		Compressor:       p.Compressor,
		Manager:          p.Manager,
		Logger:           p.Logger,
		Executor:         p.Executor,
		Store:            p.Store,
	}

	tasks = append(tasks, &task)

	for _, question := range task.ProjectStructure.Questions {
		requirements = append(requirements, &QuestionRequirement{
			Question: *question,
			Prompter: p.Prompter,
			Logger:   p.Logger,
			Executor: p.Executor,
			Manager:  p.Manager,
			Store:    p.Store,
		})
	}

	requirements = append(requirements, &templateRequirement{
		Prompter: p.Prompter,
		Store:    p.Store,
		Values:   task.ProjectStructure.Values,
	})

	return tasks, requirements, nil
}

func (p *ProjectStructureRequirement) setLanguageChecker(projectStructureData *model.ProjectStructureData) {
	checker := langs.NewLanguageChecker(projectStructureData)
	p.Store.SetLanguageChecker(checker)
}

func (p *projectStructureTask) Complete() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}

	moduleName := p.Store.GetValue(store.ModuleName)

	p.Logger.LogInfo("Extracting files...")
	if err := p.Compressor.UncompressFromUrl(p.ProjectStructure.URL); err != nil {
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
