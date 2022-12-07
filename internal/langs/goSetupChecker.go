package langs

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"os"
	"os/exec"
	"strings"
)

var (
	atSign                  = "@"
	latestVersion           = "latest"
	ErrDependencyIsNotValid = errors.New("dependency is not valid")
)

const GoCommand = "go"

type golangSetupChecker struct {
	Logger logger.Logger
	Store  store.Store
}

type CommandData struct {
	WorkingDir *string
	Command    string
	Args       []string
}

func NewGolangSetupChecker(Logger logger.Logger, Store store.Store) Checker {
	return &golangSetupChecker{
		Logger: Logger,
		Store:  Store,
	}
}

func (g *golangSetupChecker) Setup() error {
	moduleName := g.Store.GetValue(store.ModuleName)

	g.Logger.LogInfo(fmt.Sprintf("module name will be -> %s", moduleName))

	if err := g.EditGoModule(); err != nil {
		return err
	}

	g.Logger.LogInfo(fmt.Sprintf("module name was changed to -> %s", moduleName))
	return nil
}

func (g *golangSetupChecker) CheckDependency(dependency interface{}) error {
	val, ok := dependency.(string)
	if !ok {
		return fmt.Errorf("go dependecy must be string")
	}

	if len(strings.TrimSpace(val)) == 0 {
		return ErrDependencyIsNotValid
	}

	return nil
}

func (g *golangSetupChecker) CheckSetup() error {
	_, err := exec.LookPath(GoCommand)
	if err != nil {
		return fmt.Errorf("could not find %s in PATH. make sure that %s installed", "go", "go")
	}
	return nil

}

func (g *golangSetupChecker) GetDependency(dependency interface{}) error {
	val := dependency.(string)
	url := val
	version := latestVersion

	atIndex := strings.LastIndex(val, atSign)
	if atIndex != -1 {
		url = val[:atIndex]
		version = val[atIndex+1:]
	}

	goDependency := Dependency{
		Url:     &url,
		Version: &version,
	}

	g.Logger.LogInfo(fmt.Sprintf("Adding -> %s", goDependency.String()))

	data := &CommandData{
		Command: "go",
		Args:    []string{"get", goDependency.String()},
	}

	if err := g.RunCommand(data); err != nil {
		return err
	}

	g.Logger.LogInfo(fmt.Sprintf("Added  -> %s", goDependency.String()))

	return nil
}

func (g *golangSetupChecker) CleanUp() error {

	if err := g.executeFmt(); err != nil {
		return err
	}

	if err := g.executeTidy(); err != nil {
		return err
	}

	if err := g.downloadDependencies(); err != nil {
		return err
	}

	return nil
}

func (g *golangSetupChecker) executeFmt() error {
	g.Logger.LogInfo("Executing go fmt ./...")
	formatTask := &CommandData{
		Command: GoCommand,
		Args:    []string{"fmt", "./..."},
	}
	return g.RunCommand(formatTask)
}

func (g *golangSetupChecker) downloadDependencies() error {
	g.Logger.LogInfo("Executing go mod download")
	formatTask := &CommandData{
		Command: GoCommand,
		Args:    []string{"mod", "download"},
	}
	return g.RunCommand(formatTask)
}

func (g *golangSetupChecker) executeTidy() error {
	g.Logger.LogInfo("Executing go mod tidy")
	tidyTask := &CommandData{
		Command: GoCommand,
		Args:    []string{"mod", "tidy"},
	}
	return g.RunCommand(tidyTask)
}

func (g *golangSetupChecker) RunCommand(data *CommandData) error {
	if data.WorkingDir == nil {
		projectFullPath := g.Store.GetValue(store.ProjectFullPath)
		err := os.Chdir(projectFullPath)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(data.Command, data.Args...)

	err := cmd.Run()
	if err != nil {
		//	log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
}

type (
	Dependency struct {
		Url     *string
		Version *string
	}
)

func (d *Dependency) String() string {
	return fmt.Sprintf("%s@%s", *d.Url, *d.Version)
}

func (g *golangSetupChecker) EditGoModule() error {
	projectFullPath := g.Store.GetValue(store.ProjectFullPath)
	moduleName := g.Store.GetValue(store.ModuleName)

	args := make([]string, 0)

	if g.hasGoModule(projectFullPath) {
		args = append(args, "mod", "edit", "-module", moduleName)
	} else {
		args = append(args, "mod", "init", moduleName)
	}

	data := &CommandData{
		Command: "go",
		Args:    args,
	}

	return g.RunCommand(data)
}

func (g *golangSetupChecker) hasGoModule(projectDirectory string) bool {
	path := fmt.Sprintf("%s/go.mod", projectDirectory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
