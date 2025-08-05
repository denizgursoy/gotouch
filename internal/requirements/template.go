package requirements

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/denizgursoy/gotouch/internal/template"
)

const (
	ChangeValues = "Do you want to edit values?"
	EditValues   = "Press Enter to change. Values will be saved when you exit"
	YamLPattern  = "*.yaml"
)

type (
	templateRequirement struct {
		Prompter prompter.Prompter `validate:"required"`
		Store    store.Store       `validate:"required"`
		Template *template.Template
	}

	templateTask struct {
		Store    store.Store    `validate:"required"`
		Values   map[string]any `validate:"required"`
		Template *template.Template
	}
)

func (t *templateRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	tasks := make([]model.Task, 0)

	values := t.Store.GetCustomValues()

	templateTsk := &templateTask{
		Store:    t.Store,
		Values:   values,
		Template: t.Template,
	}

	if len(values) != 0 {
		yes, promptError := t.Prompter.AskForYesOrNo(ChangeValues)
		if promptError != nil {
			return nil, nil, promptError
		}

		if yes {
			marshal, marshallError := yaml.Marshal(values)
			if marshallError != nil {
				return nil, nil, marshallError
			}

			multilineString, multilineError := t.Prompter.AskForMultilineString(EditValues, string(marshal), YamLPattern)
			if multilineError != nil {
				return nil, nil, multilineError
			}

			var output map[string]any
			if unmarshallError := yaml.Unmarshal([]byte(multilineString), &output); unmarshallError != nil {
				return nil, nil, unmarshallError
			}

			templateTsk.Values = output
		}
	}
	tasks = append(tasks, templateTsk)

	return tasks, nil, nil
}

func (t *templateTask) Complete(context.Context) error {
	path := t.Store.GetValue(store.ProjectFullPath)
	t.combineWithDefaultValues()

	paths := make([]string, 0)

	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if err := t.AddSimpleTemplate(filePath); err != nil {
					prefix := strings.TrimPrefix(filePath, path)
					return fmt.Errorf("could not template file=%s, error=%v\n", prefix, err)
				}
			}
			paths = append(paths, filePath)
			return nil
		})
	if err != nil {
		return err
	}

	if err = t.templateDirectoryNames(paths); err != nil {
		return err
	}
	return err
}

func (t *templateTask) AddSimpleTemplate(path string) error {
	v, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fileWriter, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	err = t.Template.ExecuteContent(fileWriter, t.Values, v)
	if err != nil {
		return err
	}
	return nil
}

func (t *templateTask) templateDirectoryNames(folders []string) error {
	sort.Slice(folders, func(i, j int) bool {
		countI := strings.Count(folders[i], string(os.PathSeparator))
		countJ := strings.Count(folders[j], string(os.PathSeparator))
		return countI < countJ
	})

	for len(folders) > 0 {
		oldName := folders[0]
		newName, err := t.Template.Execute(t.Values, oldName)
		if err != nil {
			return err
		}

		folders = append(folders[:0], folders[0+1:]...)

		if _, err := os.Stat(newName); os.IsNotExist(err) {
			err = os.Rename(oldName, newName)
			if err != nil {
				return err
			}
			for i := range folders {
				if strings.HasPrefix(folders[i], oldName) {
					folders[i] = strings.ReplaceAll(folders[i], oldName, newName)
				}
			}
		}
	}

	return nil
}

func (t *templateTask) combineWithDefaultValues() {
	combinedValues := map[string]any{}
	if t.Values != nil {
		combinedValues = t.Values
	}

	for key, value := range t.getDefaultValues() {
		combinedValues[key] = value
	}

	for key, value := range t.Store.GetValues() {
		combinedValues[key] = value
	}
	t.Values = combinedValues
}

func (t *templateTask) getDefaultValues() map[string]any {
	defaultValues := make(map[string]any, 0)

	defaultValues[store.ProjectName] = t.Store.GetValue(store.ProjectName)
	defaultValues[store.ProjectFullPath] = t.Store.GetValue(store.ProjectFullPath)
	defaultValues[store.ModuleName] = t.Store.GetValue(store.ModuleName)
	defaultValues[store.WorkingDirectory] = t.Store.GetValue(store.WorkingDirectory)

	return defaultValues
}
