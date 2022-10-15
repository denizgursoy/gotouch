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
		Short: "Helps you create new Go Projects",
		Long:  `Tag`,
		Run:   GetCreateCommandHandler(cmdr),
	}
	createCommand.Flags().StringP(FileFlagName, "f", "", "input file")

	createCommand.AddCommand(CreatePackageCommand(cmdr))
	createCommand.AddCommand(CreateValidateCommand(cmdr))

	return createCommand
}

func CreatePackageCommand(cmdr operator.Operator) *cobra.Command {
	packageCommand := &cobra.Command{
		Use:   "package",
		Short: "createYourZip",
		Long:  `Tag`,
		Run:   GetPackageCommandHandler(cmdr),
	}

	packageCommand.Flags().StringP(SourceDirectoryFlagName, "s", ".", "source directory")
	packageCommand.Flags().StringP(TargetDirectoryFlagName, "t", "..", "target directory")

	return packageCommand
}

func CreateValidateCommand(cmdr operator.Operator) *cobra.Command {
	validateCommand := &cobra.Command{
		Use:   "validate",
		Short: "Validation Check for YAML files",
		Long:  `Tag`,
		Run:   GetValidateCommandHandler(cmdr),
	}

	validateCommand.Flags().StringP(FileFlagName, "f", "", "input file")

	return validateCommand
}
