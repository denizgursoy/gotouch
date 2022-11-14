package requirements

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"gopkg.in/yaml.v2"
)

const (
	ChangeValues = "Do you want to edit values?"
	EditValues   = "Press Enter to change. Values will be saved when you exit"
	YamLPattern  = "*.yaml"
)

type (
	templateRequirement struct {
		Prompter   prompter.Prompter `validate:"required"`
		Store      store.Store       `validate:"required"`
		Values     interface{}       `validate:"required"`
		Delimeters string
	}

	templateTask struct {
		Store      store.Store `validate:"required"`
		Values     interface{} `validate:"required"`
		Delimeters string
		Template   *template.Template
	}
)

func (t *templateRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	tasks := make([]model.Task, 0)
	templateTsk := &templateTask{
		Store:      t.Store,
		Values:     t.Values,
		Delimeters: t.Delimeters,
	}

	if t.Values != nil {
		yes, promptError := t.Prompter.AskForYesOrNo(ChangeValues)
		if promptError != nil {
			return nil, nil, promptError
		}

		if yes {
			marshal, marshallError := yaml.Marshal(t.Values)
			if marshallError != nil {
				return nil, nil, marshallError
			}

			multilineString, multilineError := t.Prompter.AskForMultilineString(EditValues, string(marshal), YamLPattern)
			if multilineError != nil {
				return nil, nil, multilineError
			}

			var output interface{}
			if unmarshallError := yaml.Unmarshal([]byte(multilineString), &output); unmarshallError != nil {
				return nil, nil, unmarshallError
			}

			templateTsk.Values = output
		}
	}
	tasks = append(tasks, templateTsk)

	return tasks, nil, nil
}

func (t *templateTask) Complete() error {
	path := t.Store.GetValue(store.ProjectFullPath)
	t.combineWithDefaultValues()

	delimiters := strings.Fields(t.Delimeters)
	templ := template.New("task")

	if len(delimiters) > 0 {
		templ.Delims(delimiters[0], delimiters[1])
	}
	t.Template = templ

	folders := make([]string, 0)

	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				t.AddSimpleTemplate(filePath)
			} else {
				folders = append(folders, filePath)
			}
			return nil
		})

	if err = t.templateDirectoryNames(folders); err != nil {
		return err
	}
	return err
}

func (t *templateTask) AddSimpleTemplate(path string) {
	files, err := t.Template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = files.Execute(f, t.Values)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *templateTask) templateDirectoryNames(folders []string) error {
	sort.Slice(folders, func(i, j int) bool {
		countI := strings.Count(folders[i], string(os.PathSeparator))
		countJ := strings.Count(folders[j], string(os.PathSeparator))
		return countI < countJ
	})

	for len(folders) > 0 {
		oldName := folders[0]
		parse, parseError := t.Template.Parse(oldName)
		if parseError != nil {
			return parseError
		}
		bufferString := bytes.NewBufferString("")
		executeError := parse.Execute(bufferString, t.Values)

		if executeError != nil {
			return parseError
		}
		newName := bufferString.String()

		folders = append(folders[:0], folders[0+1:]...)

		if _, err := os.Stat(newName); os.IsNotExist(err) {
			err = os.Rename(oldName, newName)
			if err != nil {
				return err
			}
			for i, _ := range folders {
				if strings.HasPrefix(folders[i], oldName) {
					folders[i] = strings.ReplaceAll(folders[i], oldName, newName)
				}
			}
		}
	}

	return nil
}

func (t *templateTask) combineWithDefaultValues() {
	combinedValues := map[interface{}]interface{}{}
	if t.Values != nil {
		combinedValues = t.Values.(map[interface{}]interface{})
	}

	for key, value := range t.getDefaultValues() {
		combinedValues[key] = value
	}

	for key, value := range t.Store.GetValues() {
		combinedValues[key] = value
	}
	t.Values = combinedValues
}

func (t *templateTask) getDefaultValues() map[interface{}]interface{} {
	defaultValues := make(map[interface{}]interface{}, 0)

	defaultValues[store.ProjectName] = t.Store.GetValue(store.ProjectName)
	defaultValues[store.ProjectFullPath] = t.Store.GetValue(store.ProjectFullPath)
	defaultValues[store.ModuleName] = t.Store.GetValue(store.ModuleName)
	defaultValues[store.WorkingDirectory] = t.Store.GetValue(store.WorkingDirectory)

	return defaultValues
}
