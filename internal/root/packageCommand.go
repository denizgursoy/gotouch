package root

import (
	"github.com/denizgursoy/gotouch/internal/commander"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/spf13/cobra"
	"strings"
)

const (
	SourceDirectoryFlagName = "source"
	TargetDirectoryFlagName = "target"
)

func GetPackageCommandHandler(cmdr commander.Commander) CommandHandler {
	return func(cmd *cobra.Command, args []string) {
		lgr := logger.NewLogger()

		flags := cmd.Flags()
		sourceDirectoryPath, inputError := flags.GetString(SourceDirectoryFlagName)
		targetDirectoryPath, inputError := flags.GetString(TargetDirectoryFlagName)

		if inputError != nil {
			lgr.LogErrorIfExists(inputError)
			return
		}

		sourcePointer := &sourceDirectoryPath
		if len(strings.TrimSpace(sourceDirectoryPath)) == 0 {
			sourcePointer = nil
		}

		targetPointer := &targetDirectoryPath
		if len(strings.TrimSpace(targetDirectoryPath)) == 0 {
			sourcePointer = nil
		}

		options := commander.PackageCommandOptions{
			Compressor:      compressor.GetInstance(),
			Logger:          lgr,
			SourceDirectory: sourcePointer,
			TargetDirectory: targetPointer,
		}
		err := cmdr.CompressDirectory(&options)
		lgr.LogErrorIfExists(err)
	}
}
