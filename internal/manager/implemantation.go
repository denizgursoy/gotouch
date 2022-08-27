package manager

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/executor"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	urls        []string
	index       = 0
	Environment = "prod"
)

type (
	fManager struct {
		Executor executor.Executor
	}

	Dependency struct {
		Url     *string
		Version *string
	}
)

func newFileManager() Manager {
	return &fManager{
		Executor: executor.GetInstance(),
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

func (f *fManager) AddDependency(dependency string) error {
	panic("implement me")
}

func (d *Dependency) String() string {
	return fmt.Sprintf("%s@%s", *d.Url, *d.Version)
}

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
			split := strings.Split(line, " ")
			ints := make([]byte, 0)

			for _, s := range split {
				atoi, _ := strconv.Atoi(s)
				ints = append(ints, byte(atoi))
			}
			urls = append(urls, string(ints))
		}
	}
}

func (f *fManager) EditGoModule(projectName, folderName string) error {
	workingDirectory := f.GetExtractLocation()
	projectDirectory := fmt.Sprintf("%s/%s", workingDirectory, folderName)

	if err := os.Chdir(projectDirectory); err != nil {
		return err
	}

	args := make([]string, 0)

	if f.hasGoModule(projectDirectory) {
		args = append(args, "mod", "edit", "-module", projectName)
	} else {
		args = append(args, "mod", "init", projectName)
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
