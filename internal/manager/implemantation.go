package manager

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
)

type (
	fManager struct {
		Executor executor.Executor `validate:"required"`
		Store    store.Store       `validate:"required"`
		Logger   logger.Logger     `validate:"required"`
	}
)

func init() {
	manager = newFileManager()
}

func newFileManager() Manager {
	return &fManager{
		Executor: executor.GetInstance(),
		Store:    store.GetInstance(),
		Logger:   logger.NewLogger(),
	}
}

func (f *fManager) CreateDirectoryIfNotExists(directoryName string) error {
	return os.Mkdir(directoryName, os.ModePerm)
}

func (f *fManager) GetExtractLocation() string {
	return GetExtractLocation()
}

func (f *fManager) EditGoModule() error {
	projectFullPath := f.Store.GetValue(store.ProjectFullPath)
	moduleName := f.Store.GetValue(store.ModuleName)

	args := make([]string, 0)

	if f.hasGoModule(projectFullPath) {
		args = append(args, "mod", "edit", "-module", moduleName)
	} else {
		args = append(args, "mod", "init", moduleName)
	}

	data := &executor.CommandData{
		Command: "go",
		Args:    args,
	}

	return f.Executor.RunCommand(data)
}

func (f *fManager) hasGoModule(projectDirectory string) bool {
	path := fmt.Sprintf("%s/go.mod", projectDirectory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (f *fManager) CreateFile(reader io.ReadCloser, path string) error {
	fullPath := filepath.Join(f.Store.GetValue(store.ProjectFullPath), path)
	directoryOfFile := filepath.Dir(fullPath)

	err := os.MkdirAll(directoryOfFile, os.ModePerm)
	if err != nil {
		return err
	}
	f.Logger.LogInfo(fmt.Sprintf("Creating file -> %s", fullPath))

	create, createError := os.Create(fullPath)
	if createError != nil {
		return createError
	}
	defer create.Close()

	all, readError := ioutil.ReadAll(reader)
	if readError != nil {
		return readError
	}
	_, writeError := create.Write(all)
	if writeError != nil {
		return writeError
	}

	f.Logger.LogInfo(fmt.Sprintf("Created file  -> %s", fullPath))
	return nil
}
