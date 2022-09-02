package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/skratchdot/open-golang/open"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	EnterValues = "Do you want to edit values?"
	Ready       = "Is file ready?"
)

type (
	templateRequirement struct {
		Prompter prompter.Prompter `validate:"required"`
		Store    store.Store       `validate:"required"`
		Values   interface{}       `validate:"required"`
	}

	templateTask struct {
		Store  store.Store
		Values interface{} `validate:"required"`
	}
)

func (t *templateRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	if t.Values != nil {
		yes, err2 := t.Prompter.AskForYesOrNo(EnterValues)
		if err2 != nil {
			return nil, nil, err2
		}
		tasks := make([]model.Task, 0)

		if yes {
			marshal, err2 := yaml.Marshal(t.Values)
			if err2 != nil {
				return nil, nil, err2
			}

			temp, err2 := os.CreateTemp("", "*.yaml")
			if err2 != nil {
				return nil, nil, err2
			}

			_, err2 = temp.Write(marshal)
			if err2 != nil {
				return nil, nil, err2
			}

			err2 = open.Run(temp.Name())
			defer func() {
				err2 := os.Remove(temp.Name())
				if err2 != nil {
					log.Fatal(err2)
				}
			}()

			if err2 != nil {
				return nil, nil, err2
			}

			yes, err2 := t.Prompter.AskForYesOrNo(fmt.Sprintf("%s (%s)", Ready, temp.Name()))
			if yes == false || err2 != nil {
				return nil, nil, err2
			}

			all, err2 := ioutil.ReadFile(temp.Name())
			var output interface{}
			err2 = yaml.Unmarshal(all, &output)
			templateTsk := &templateTask{
				Store:  t.Store,
				Values: output,
			}
			tasks = append(tasks, templateTsk)
			return tasks, nil, nil
		}

	}

	return nil, nil, nil
}

func (t *templateTask) Complete() error {
	path := t.Store.GetValue(store.ProjectFullPath)
	t.combineWithDefaultValues()

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() != ".keep" {
				t.AddSimpleTemplate(path)
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (t *templateTask) AddSimpleTemplate(path string) {
	files, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = files.Execute(f, t.Values)
}

func (t *templateTask) combineWithDefaultValues() {
	m := t.Values.(map[interface{}]interface{})
	for key, value := range t.getDefaultValues() {
		m[key] = value
	}
	t.Values = m
}

func (t *templateTask) getDefaultValues() map[interface{}]interface{} {
	defaultValues := make(map[interface{}]interface{}, 0)

	defaultValues[store.ProjectName] = t.Store.GetValue(store.ProjectName)
	defaultValues[store.ProjectFullPath] = t.Store.GetValue(store.ProjectFullPath)
	defaultValues[store.ModuleName] = t.Store.GetValue(store.ModuleName)
	defaultValues[store.WorkingDirectory] = t.Store.GetValue(store.WorkingDirectory)

	return defaultValues
}
