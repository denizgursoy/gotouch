package root

import (
	"github.com/denizgursoy/gotouch/internal/langs"
	"testing"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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

			expectedCall := &operator.CreateNewProjectOptions{
				Lister:          lister.GetInstance(),
				Prompter:        prompter.GetInstance(),
				Manager:         manager.GetInstance(),
				Compressor:      compressor.GetInstance(),
				Executor:        executor.GetInstance(),
				Logger:          logger.NewLogger(),
				Path:            argument.pointer,
				Store:           store.GetInstance(),
				LanguageChecker: langs.GetInstance(),
			}

			mockCommander.EXPECT().CreateNewProject(gomock.Eq(expectedCall))

			command := CreateRootCommand(mockCommander)
			command.SetArgs(getCreateTestArguments(argument.flag))

			err := command.Execute()
			require.Nil(t, err)
		}
	})
}

func getCreateTestArguments(filePath string) []string {
	return []string{"-f", filePath}
}
