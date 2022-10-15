package root

import (
	"strings"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/spf13/cobra"
)

const (
	SourceDirectoryFlagName = "source"
	TargetDirectoryFlagName = "target"
)

func GetPackageCommandHandler(cmdr operator.Operator) CommandHandler {
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

		options := operator.CompressDirectoryOptions{
			Compressor:      compressor.GetInstance(),
			Logger:          lgr,
			SourceDirectory: sourcePointer,
			TargetDirectory: targetPointer,
		}
		err := cmdr.CompressDirectory(&options)
		lgr.LogErrorIfExists(err)
	}
}
