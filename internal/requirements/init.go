package requirements

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	initRequirement struct {
		Store  store.Store
		Logger logger.Logger
	}
	initTask struct {
		Store  store.Store
		Logger logger.Logger
	}
)

func (i *initRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	tasks := make([]model.Task, 0)

	tasks = append(tasks, &initTask{
		Store:  i.Store,
		Logger: i.Logger,
	})
	return tasks, nil, nil
}

func (i *initTask) Complete() error {
	if err := validator.New().Struct(i); err != nil {
		return err
	}

	projectFullPath := i.Store.GetValue(store.ProjectFullPath)
	initFileAddress := filepath.Join(projectFullPath, InitFileName)

	_, err := os.Stat(initFileAddress)
	if err == nil {
		i.Logger.LogInfo("Executing " + InitFileName)
		defer os.Remove(initFileAddress)
		err := os.Chmod(initFileAddress, 0777)
		if err != nil {
			return err
		}
		err = executeInitFile(i.Store)
		if err != nil {
			return err
		}
		i.Logger.LogInfo("Executed " + InitFileName)
		return nil
	}
	return nil
}

type CommandData struct {
	WorkingDir *string
	Command    string
	Args       []string
}

func RunCommand(data *CommandData, str store.Store) error {
	if data.WorkingDir == nil {
		projectFullPath := str.GetValue(store.ProjectFullPath)
		err := os.Chdir(projectFullPath)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(data.Command, data.Args...)

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		//	log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}
