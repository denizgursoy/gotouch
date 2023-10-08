package commands

import (
	"strconv"
	"strings"

	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/config"

	"github.com/denizgursoy/gotouch/internal/operator"

	"github.com/spf13/cobra"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
)

const (
	FileFlagName   = "file"
	InlineFlagName = "inline"
)

type (
	CommandHandler func(cmd *cobra.Command, args []string)
)

func GetCreateCommandHandler(cmdr operator.Operator) CommandHandler {
	return func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		newLogger := logger.NewLogger()

		filePath, inputError := flags.GetString(FileFlagName)
		if inputError != nil {
			newLogger.LogErrorIfExists(inputError)
			return
		}

		pointer := &filePath
		if len(strings.TrimSpace(filePath)) == 0 {
			pointer = nil
		}

		inline, inlineError := flags.GetBool(InlineFlagName)
		if inlineError != nil {
			newLogger.LogErrorIfExists(inlineError)
			return
		}

		appStore := store.GetInstance()
		appStore.SetValue(store.Inline, strconv.FormatBool(inline))

		options := operator.CreateNewProjectOptions{
			Lister:        lister.GetInstance(),
			Prompter:      prompter.GetInstance(),
			Manager:       manager.GetInstance(),
			Compressor:    compressor.GetInstance(),
			Executor:      executor.GetInstance(),
			Logger:        newLogger,
			Path:          pointer,
			Store:         appStore,
			Cloner:        cloner.GetInstance(),
			CommandRunner: commandrunner.GetInstance(appStore),
			ConfigManager: config.NewConfigManager(newLogger),
		}

		err := cmdr.CreateNewProject(&options)
		options.Logger.LogErrorIfExists(err)
	}
}
