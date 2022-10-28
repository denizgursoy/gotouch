package commands

import (
	"log"
	"strings"

	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/spf13/cobra"
)

func GetValidateCommandHandler(cmdr operator.Operator) CommandHandler {
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

		options := operator.ValidateYamlOptions{
			Lister: lister.GetInstance(),
			Logger: logger.NewLogger(),
			Path:   point,
		}
		err := cmdr.ValidateYaml(&options)
		options.Logger.LogErrorIfExists(err)
	}
}
