package requirements

import (
	"context"
	"fmt"
	"strings"

	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/template"

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
		VCSDetector     cloner.VCSDetector   `validate:"required"`
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
		VCSDetector      cloner.VCSDetector          `validate:"required"`
		Client           model.HttpRequester
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

	// TODO: test
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
	tasks = append(tasks, p.getProjectStructureTask(selectedPS))
	requirements = append(requirements, nameRequirements...)

	resourceTasks := getTasks(selectedPS.Resources, p.Logger, p.Manager, p.LanguageChecker, p.Store)
	tasks = append(tasks, resourceTasks...)
	tasks = append(tasks, p.getStartUpTask())

	p.Store.AddCustomValues(selectedPS.CustomValues)
	p.Store.AddValues(selectedPS.Values)

	for _, question := range selectedPS.Questions {
		requirements = append(requirements, p.getQuestionRequirement(question))
	}

	requirements = append(requirements, p.getTemplateRequirement(selectedPS))
	requirements = append(requirements, p.getInitRequirement())
	requirements = append(requirements, p.getCleanupRequirement())

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

func (p *ProjectStructureRequirement) getStartUpTask() *startupTask {
	return &startupTask{
		Store:  p.Store,
		Logger: p.Logger,
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
	if len(p.ProjectsData) == 1 {
		return p.ProjectsData[0], nil
	}

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

func (p *projectStructureTask) Complete(ctx context.Context) error {
	if err := validator.New().StructCtx(ctx, p); err != nil {
		return err
	}

	url := p.ProjectStructure.URL

	if len(strings.TrimSpace(url)) != 0 {
		vcs, err := p.VCSDetector.DetectVCS(ctx, p.Client, url)
		if err != nil {
			return err
		}

		switch vcs {
		case cloner.VCSGit:
			if err := p.Cloner.CloneFromUrl(ctx, url, p.ProjectStructure.Branch); err != nil {
				return err
			}
		default:
			if err := p.Compressor.UncompressFromUrl(ctx, url); err != nil {
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
