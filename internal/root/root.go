package root

import (
	"os"

	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCommand := CreateRootCommand(operator.GetInstance())
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func CreateRootCommand(cmdr operator.Operator) *cobra.Command {
	createCommand := &cobra.Command{
		Use:   "gotouch",
		Short: "Helps you create new  projects",
		Long: `Parses properties yaml provided with file flag. If no flag is provided, Gotouch will use the yaml at
https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml. 
Gotouch will list the project structures in the yaml to user in order to make selection.
Gotouch will ask for project name. Project name will be used to create directory to which archive file will be extracted.
If selected project structure's language is go, project name will be used as module name. See below:

Project Name						Module Name 						Directory Name
my-app								my-app								my-app
github.com/my-account/my-app		github.com/my-account/my-app		my-app

If no go.mod file exists in the template, Gotouch will create one with the module name.
Gotouch will ask the questions in the selected project structure in order.`,
		Run: GetCreateCommandHandler(cmdr),
	}
	createCommand.Flags().StringP(FileFlagName, "f", "", "input properties yaml")

	createCommand.AddCommand(CreatePackageCommand(cmdr))
	createCommand.AddCommand(CreateValidateCommand(cmdr))

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
