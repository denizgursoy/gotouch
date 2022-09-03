package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	SourceDirectoryFlagName = "source"
	TargetDirectoryFlagName = "target"
)

var (
	packageCommand = &cobra.Command{
		Use:   "package",
		Short: "createYourZip",
		Long:  `Tag`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			sourceDirectoryPath, inputError := flags.GetString(SourceDirectoryFlagName)
			targetDirectoryPath, inputError := flags.GetString(TargetDirectoryFlagName)

			lgr := logger.NewLogger()
			lgr.LogErrorIfExists(inputError)

			if inputError != nil {
				os.Exit(1)
			}

			sourcePointer := &sourceDirectoryPath
			if len(strings.TrimSpace(sourceDirectoryPath)) == 0 {
				sourcePointer = nil
			}

			targetPointer := &targetDirectoryPath
			if len(strings.TrimSpace(targetDirectoryPath)) == 0 {
				sourcePointer = nil
			}

			options := PackageCommandOptions{
				Compressor:      compressor.GetInstance(),
				Logger:          lgr,
				SourceDirectory: sourcePointer,
				TargetDirectory: targetPointer,
			}
			err := CompressDirectory(&options)
			lgr.LogErrorIfExists(err)
		},
	}
)

type (
	PackageCommandOptions struct {
		SourceDirectory *string               `validate:"required,dir"`
		TargetDirectory *string               `validate:"omitempty,dir"`
		Compressor      compressor.Compressor `validate:"required"`
		Logger          logger.Logger
	}
)

func init() {
	packageCommand.Flags().StringP(SourceDirectoryFlagName, "s", ".", "source directory")
	packageCommand.Flags().StringP(TargetDirectoryFlagName, "t", ".", "target directory")
}

func CompressDirectory(opts *PackageCommandOptions) error {
	targetDirectory := ""
	if opts.TargetDirectory != nil {
		targetDirectory = *opts.TargetDirectory
	}
	return opts.Compressor.CompressDirectory(*opts.SourceDirectory, targetDirectory)
}
