package requirements

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/template"
	"strings"

	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/langs"
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
		Cloner          cloner.Cloner        `validate:"required"`
		CommandRunner   commandrunner.Runner `validate:"required"`
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

	selectedPS, err := p.SelectProject()
	if err != nil {
		return nil, nil, err
	}

	//TODO: test
	p.LanguageChecker = langs.GetChecker(selectedPS.Language, p.Logger, p.Store, p.CommandRunner)
	if setupError := p.LanguageChecker.CheckSetup(); setupError != nil {
		return nil, nil, setupError
	}

	nameRequirement := p.getNameRequirement()

	nameTasks, nameRequirements, err := nameRequirement.AskForInput()
	if err != nil {
		return nil, nil, err
	}

	tasks := make([]model.Task, 0)
	requirements := make([]model.Requirement, 0)

	tasks = append(tasks, nameTasks...)
	requirements = append(requirements, nameRequirements...)

	resourceTasks := getTasks(selectedPS.Resources, p.Logger, p.Manager, p.LanguageChecker, p.Store)

	if resourceTasks != nil {
		tasks = append(tasks, resourceTasks...)
	}

	tasks = append(tasks, p.getProjectStructureTask(selectedPS))

	for _, question := range selectedPS.Questions {
		requirements = append(requirements, p.getQuestionRequirement(question))
	}

	requirements = append(requirements, p.getTemplateRequirement(selectedPS))
	requirements = append(requirements, p.getCleanupRequirement())
	requirements = append(requirements, p.getInitRequirement())

	return tasks, requirements, nil
}

func (p *ProjectStructureRequirement) getNameRequirement() *ProjectNameRequirement {
	return &ProjectNameRequirement{
		Prompter: p.Prompter,
		Manager:  p.Manager,
		Logger:   p.Logger,
		Store:    p.Store,
	}
}

func (p *ProjectStructureRequirement) getProjectStructureTask(selectedPS *model.ProjectStructureData) *projectStructureTask {
	return &projectStructureTask{
		ProjectStructure: selectedPS,
		Compressor:       p.Compressor,
		Manager:          p.Manager,
		Logger:           p.Logger,
		Executor:         p.Executor,
		Store:            p.Store,
		LanguageChecker:  p.LanguageChecker,
		Cloner:           p.Cloner,
	}
}

func (p *ProjectStructureRequirement) getQuestionRequirement(question *model.Question) *QuestionRequirement {
	return &QuestionRequirement{
		Question:        *question,
		Prompter:        p.Prompter,
		Logger:          p.Logger,
		Executor:        p.Executor,
		Manager:         p.Manager,
		Store:           p.Store,
		LanguageChecker: p.LanguageChecker,
	}
}

func (p *ProjectStructureRequirement) getTemplateRequirement(selectedPS *model.ProjectStructureData) *templateRequirement {
	templateWithDelimiters := GetTemplate(selectedPS)

	return &templateRequirement{
		Prompter: p.Prompter,
		Store:    p.Store,
		Values:   selectedPS.Values,
		Template: templateWithDelimiters,
	}
}

func (p *ProjectStructureRequirement) getCleanupRequirement() *cleanupRequirement {
	return &cleanupRequirement{
		LanguageChecker: p.LanguageChecker,
	}
}

func (p *ProjectStructureRequirement) getInitRequirement() *initRequirement {
	return &initRequirement{
		Store:         p.Store,
		Logger:        p.Logger,
		CommandRunner: p.CommandRunner,
	}
}

func (p *ProjectStructureRequirement) SelectProject() (*model.ProjectStructureData, error) {
	options := make([]fmt.Stringer, 0)
	for _, datum := range p.ProjectsData {
		options = append(options, datum)
	}

	selected, err := p.Prompter.AskForSelectionFromList(SelectProjectTypeDirection, options)
	if err != nil {
		return nil, err
	}

	projectStructureData := selected.(*model.ProjectStructureData)
	return projectStructureData, nil
}

func (p *projectStructureTask) Complete() error {
	if err := validator.New().Struct(p); err != nil {
		return err
	}

	url := p.ProjectStructure.URL

	if len(strings.TrimSpace(url)) != 0 {
		if strings.HasSuffix(url, ".git") {
			if err := p.Cloner.CloneFromUrl(url, p.ProjectStructure.Branch); err != nil {
				return err
			}
		} else {
			if err := p.Compressor.UncompressFromUrl(url); err != nil {
				return err
			}
		}
	}

	if preTaskError := p.LanguageChecker.Setup(); preTaskError != nil {
		return preTaskError
	}

	return nil
}

func GetTemplate(p *model.ProjectStructureData) *template.Template {
	t := template.New()

	t.SetSprigFuncs()

	delimiters := strings.Fields(p.Delimiters)
	if len(delimiters) > 0 {
		t.SetDelims(delimiters[0], delimiters[1])
	}

	return t
}
