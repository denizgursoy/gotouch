package root

import (
	"github.com/denizgursoy/gotouch/internal/commander"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func GetValidateCommandHandler(cmdr commander.Commander) CommandHandler {
	return func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		filePath, inputError := flags.GetString(FileFlagName)
		if inputError != nil {
			log.Fatalln(inputError)
		}

		point := &filePath
		if len(strings.TrimSpace(filePath)) == 0 {
			point = nil
		}

		options := commander.ValidateCommandOptions{
			Lister: lister.GetInstance(),
			Logger: logger.NewLogger(),
			Path:   point,
		}
		err := cmdr.ValidateYaml(&options)
		options.Logger.LogErrorIfExists(err)
	}
}
