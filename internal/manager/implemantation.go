package manager

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
)

var (
	urls        []string
	index       = 0
	Environment = "prod"
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
	if manager.IsTest() {
		exPath := fmt.Sprintf("%s/input.txt", manager.GetExtractLocation())
		file, err := os.ReadFile(exPath)
		if err != nil {
			log.Println("deniz", err)
		}
		urls = make([]string, 0)
		for _, line := range strings.Split(string(file), "\n") {
			urls = append(urls, line)
		}
	}
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

func (f *fManager) GetStream() (ioReader io.ReadCloser) {
	if f.IsTest() {
		ioReader = io.NopCloser(strings.NewReader(urls[index]))
	} else {
		ioReader = os.Stdin
	}
	index++
	return
}

func (f *fManager) IsTest() bool {
	return Environment == "test"
}

func (f *fManager) GetExtractLocation() string {
	if f.IsTest() {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal("could not fetch executable information", err)
		}
		return filepath.Dir(ex)
	} else {
		return f.GetWd()
	}
}

func (f *fManager) GetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get working directory", err)
	}

	return wd
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

	err2 := os.MkdirAll(directoryOfFile, os.ModePerm)
	if err2 != nil {
		return err2
	}
	f.Logger.LogInfo(fmt.Sprintf("Creating file -> %s", fullPath))

	create, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer create.Close()

	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	_, err = create.Write(all)
	if err != nil {
		return err
	}

	f.Logger.LogInfo(fmt.Sprintf("Created file  -> %s", fullPath))
	return nil
}
