package commands

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/config"
	"go.uber.org/mock/gomock"

	"github.com/denizgursoy/gotouch/internal/cloner"

	"github.com/stretchr/testify/require"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
)

func TestGetCreateCommandHandler(t *testing.T) {
	t.Run("should create successfully", func(t *testing.T) {
		type arg struct {
			flag    string
			pointer *string
		}

		flag := "./test-input.yaml"
		arguments := []arg{
			{
				flag:    flag,
				pointer: &flag,
			},
			{
				flag:    "",
				pointer: nil,
			},
		}

		for _, argument := range arguments {
			controller := gomock.NewController(t)
			mockCommander := operator.NewMockOperator(controller)

			appStore := store.GetInstance()
			expectedCall := &operator.CreateNewProjectOptions{
				Lister:        lister.GetInstance(),
				Prompter:      prompter.GetInstance(),
				Manager:       manager.GetInstance(),
				Compressor:    compressor.GetInstance(),
				Executor:      executor.GetInstance(),
				Logger:        logger.NewLogger(),
				Path:          argument.pointer,
				Store:         appStore,
				Cloner:        cloner.GetInstance(),
				VCSDetector:   cloner.NewDefaultVCSDetector(),
				CommandRunner: commandrunner.GetInstance(appStore),
				ConfigManager: config.NewConfigManager(logger.NewLogger()),
			}

			mockCommander.EXPECT().CreateNewProject(gomock.Any(), gomock.Eq(expectedCall))

			command := CreateRootCommand(mockCommander, BuildInfo{})
			command.SetArgs(getCreateTestArguments(argument.flag))

			err := command.Execute()
			require.Nil(t, err)
		}
	})
}

func getCreateTestArguments(filePath string) []string {
	return []string{"-f", filePath}
}
