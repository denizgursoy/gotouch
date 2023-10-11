package requirements

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"

	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/store"
)

type (
	initRequirement struct {
		Store         store.Store
		Logger        logger.Logger
		CommandRunner commandrunner.Runner `validate:"required"`
	}
	initTask struct {
		Store         store.Store
		Logger        logger.Logger
		CommandRunner commandrunner.Runner `validate:"required"`
	}
)

const (
	LinuxInitFile   = "init.sh"
	WindowsInitFile = "init.bat"
)

var InitFiles = []string{LinuxInitFile, WindowsInitFile}

func (i *initRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	tasks := make([]model.Task, 0)

	tasks = append(tasks, &initTask{
		Store:         i.Store,
		Logger:        i.Logger,
		CommandRunner: i.CommandRunner,
	})
	return tasks, nil, nil
}

func (i *initTask) Complete(ctx context.Context) error {
	if err := validator.New().StructCtx(ctx, i); err != nil {
		return err
	}

	projectFullPath := i.Store.GetValue(store.ProjectFullPath)

	defer deleteInitFiles(projectFullPath)

	initFile := LinuxInitFile

	if runtime.GOOS == "windows" {
		initFile = WindowsInitFile
	}

	initFileAddress := filepath.Join(projectFullPath, initFile)

	_, err := os.Stat(initFileAddress)
	if err == nil {
		i.Logger.LogInfo("Executing " + initFile)
		if err := os.Chmod(initFileAddress, 0o777); err != nil {
			return err
		}
		if err = i.CommandRunner.Run(getCommand()); err != nil {
			return err
		}

		i.Logger.LogInfo("Executed " + initFile)
	}
	return nil
}

func deleteInitFiles(projectFullPath string) {
	for i := range InitFiles {
		initFileAddress := filepath.Join(projectFullPath, InitFiles[i])
		err := os.Remove(initFileAddress)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("could not delete file %s \n", initFileAddress)
		}
	}
}
