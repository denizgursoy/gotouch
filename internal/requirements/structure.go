package requirements

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/langs"
	"strings"

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
		ProjectsData    []*model.ProjectStructureData
		Prompter        prompter.Prompter     `validate:"required"`
		Compressor      compressor.Compressor `validate:"required"`
		Manager         manager.Manager       `validate:"required"`
		Logger          logger.Logger         `validate:"required"`
		Executor        executor.Executor     `validate:"required"`
		Store           store.Store           `validate:"required"`
		LanguageChecker langs.Checker
		Cloner          cloner.Cloner `validate:"required"`
	}

	projectStructureTask struct {
		ProjectStructure *model.ProjectStructureData `validate:"required"`
		Compressor       compressor.Compressor       `validate:"required"`
		Manager          manager.Manager             `validate:"required"`
		Executor         executor.Executor           `validate:"required"`
		Logger           logger.Logger               `validate:"required"`
		Store            store.Store                 `validate:"required"`
		LanguageChecker  langs.Checker               `validate:"required"`
		Cloner           cloner.Cloner               `validate:"required"`
	}
)

const (
	SelectProjectTypeDirection = "Select Project Type :"
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

	projectStructureData := selected.(*model.ProjectStructureData)

	//TODO: test
	p.LanguageChecker = langs.GetChecker(projectStructureData.Language, p.Logger, p.Store)
	if setupError := p.LanguageChecker.CheckSetup(); setupError != nil {
		return nil, nil, setupError
	}

	nameRequirement := &ProjectNameRequirement{
		Prompter: p.Prompter,
		Manager:  p.Manager,
		Logger:   p.Logger,
		Store:    p.Store,
	}

	nameTasks, nameRequirements, err := nameRequirement.AskForInput()
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

	task := projectStructureTask{
		ProjectStructure: projectStructureData,
		Compressor:       p.Compressor,
		Manager:          p.Manager,
		Logger:           p.Logger,
		Executor:         p.Executor,
		Store:            p.Store,
		LanguageChecker:  p.LanguageChecker,
		Cloner:           p.Cloner,
	}

	tasks = append(tasks, &task)

	for _, question := range task.ProjectStructure.Questions {
		requirements = append(requirements, &QuestionRequirement{
			Question:        *question,
			Prompter:        p.Prompter,
			Logger:          p.Logger,
			Executor:        p.Executor,
			Manager:         p.Manager,
			Store:           p.Store,
			LanguageChecker: p.LanguageChecker,
		})
	}

	requirements = append(requirements, &templateRequirement{
		Prompter: p.Prompter,
		Store:    p.Store,
		Values:   task.ProjectStructure.Values,
	})

	tasks = append(tasks, &cleanupTask{
		LanguageChecker: p.LanguageChecker,
	})

	return tasks, requirements, nil
}

func (p *projectStructureTask) Complete() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}

	url := p.ProjectStructure.URL

	if strings.HasSuffix(url, ".git") {
		if err := p.Cloner.CloneFromUrl(url); err != nil {
			return err
		}
	} else {
		if err := p.Compressor.UncompressFromUrl(url); err != nil {
			return err
		}
	}

	if preTaskError := p.LanguageChecker.Setup(); preTaskError != nil {
		return preTaskError
	}

	return nil
}
