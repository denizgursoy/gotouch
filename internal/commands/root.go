package commands

import (
	"fmt"
	"os"

	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Version     string
	BuildCommit string
	BuildDate   string
}

func Execute(info BuildInfo) {
	rootCommand := CreateRootCommand(operator.GetInstance(), info)
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func CreateRootCommand(cmdr operator.Operator, info BuildInfo) *cobra.Command {
	createCommand := &cobra.Command{
		Use:     "gotouch",
		Version: info.Version,
		Short:   "Helps you create new projects",
		Long: `Gotouch helps you create new projects from templates.
Version: ` + info.Version + ` Commit: ` + info.BuildCommit + ` Date: ` + info.BuildDate + `

Parses properties yaml provided with file flag. If no flag is provided, Gotouch will use the yaml at
https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml.

Gotouch will list the project structures in the yaml to user in order to make selection.
Gotouch will ask for project name. Project name will be used to create directory to which archive file will be extracted.
If selected project structure's language is go, project name will be used as module name. See below:

Project Name                 | Module Name                  | Directory Name
my-app                       | my-app                       | my-app
github.com/my-account/my-app | github.com/my-account/my-app | my-app

If no go.mod file exists in the template, Gotouch will create one with the module name.
Gotouch will ask the questions in the selected project structure in order.`,
		Run: GetCreateCommandHandler(cmdr),
	}
	createCommand.Flags().StringP(FileFlagName, "f", "", "input properties yaml")

	createCommand.AddCommand(CreatePackageCommand(cmdr))
	createCommand.AddCommand(CreateValidateCommand(cmdr))
	createCommand.AddCommand(CreateConfigCommand())

	return createCommand
}

func CreatePackageCommand(cmdr operator.Operator) *cobra.Command {
	packageCommand := &cobra.Command{
		Use:   "package",
		Short: "Create archive file from directory compatible with gotouch",
		Long: `Creates a tar file with gzip compression. Ignores following files/directories:

__MACOS
.DS_Store
.idea
.vscode
.git

If no flag is not set, package command creates archive file in the parent directory 
with .tar.gz extension from the files/directories in the working dir.`,
		Run: GetPackageCommandHandler(cmdr),
	}

	packageCommand.Flags().StringP(SourceDirectoryFlagName, "s", ".", "source directory")
	packageCommand.Flags().StringP(TargetDirectoryFlagName, "t", "..", "target directory")

	return packageCommand
}

func CreateValidateCommand(cmdr operator.Operator) *cobra.Command {
	validateCommand := &cobra.Command{
		Use:   "validate",
		Short: "Validation Check for YAML files",
		Long: `Checks if properties yaml can be used by gotouch.
See https://raw.githubusercontent.com/denizgursoy/gotouch/main/examples/complete-choice-example.yaml as a guide.`,
		Run: GetValidateCommandHandler(cmdr),
	}

	validateCommand.Flags().StringP(FileFlagName, "f", "", "input properties yaml")

	return validateCommand
}

func CreateConfigCommand() *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Set/Unset values",
	}

	configCommand.AddCommand(&cobra.Command{
		Use:  "set",
		Args: cobra.MatchAll(isConfigurable, cobra.ExactArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConfig(args[0], args[1])
		},
	})
	configCommand.AddCommand(&cobra.Command{
		Use:  "unset",
		Args: cobra.MatchAll(isConfigurable, cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unsetConfig(args[0])
		},
	})
	return configCommand
}

func isConfigurable(cmd *cobra.Command, args []string) error {
	configurableSettings := []string{"url"}
	for _, confArg := range configurableSettings {
		if confArg == args[0] {
			return nil
		}
	}
	return fmt.Errorf("%s is not a valid argumet", args[0])
}
