package langs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
)

var (
	atSign                  = "@"
	latestVersion           = "latest"
	ErrDependencyIsNotValid = errors.New("dependency is not valid")
)

const GoCommand = "go"

type golangSetupChecker struct {
	Logger        logger.Logger
	Store         store.Store
	CommandRunner commandrunner.Runner
}

func NewGolangSetupChecker(Logger logger.Logger, Store store.Store, runner commandrunner.Runner) Checker {
	return &golangSetupChecker{
		Logger:        Logger,
		Store:         Store,
		CommandRunner: runner,
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

func (g *golangSetupChecker) CheckDependency(dependency any) error {
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

func (g *golangSetupChecker) GetDependency(dependency any) error {
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

	data := &commandrunner.CommandData{
		Command: "go",
		Args:    []string{"get", goDependency.String()},
	}

	if err := g.CommandRunner.Run(data); err != nil {
		return err
	}

	g.Logger.LogInfo(fmt.Sprintf("Added  -> %s", goDependency.String()))

	return nil
}

func (g *golangSetupChecker) CleanUp() error {
	return nil
}

func (g *golangSetupChecker) executeFmt() error {
	g.Logger.LogInfo("Executing go fmt ./...")
	formatTask := &commandrunner.CommandData{
		Command: GoCommand,
		Args:    []string{"fmt", "./..."},
	}
	return g.CommandRunner.Run(formatTask)
}

func (g *golangSetupChecker) downloadDependencies() error {
	g.Logger.LogInfo("Executing go mod download")
	formatTask := &commandrunner.CommandData{
		Command: GoCommand,
		Args:    []string{"mod", "download"},
	}
	return g.CommandRunner.Run(formatTask)
}

func (g *golangSetupChecker) executeTidy() error {
	g.Logger.LogInfo("Executing go mod tidy")
	tidyTask := &commandrunner.CommandData{
		Command: GoCommand,
		Args:    []string{"mod", "tidy"},
	}
	return g.CommandRunner.Run(tidyTask)
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

	data := &commandrunner.CommandData{
		Command: GoCommand,
		Args:    args,
	}

	return g.CommandRunner.Run(data)
}

func (g *golangSetupChecker) hasGoModule(projectDirectory string) bool {
	path := fmt.Sprintf("%s/go.mod", projectDirectory)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
