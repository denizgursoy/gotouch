package commands

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetPackageCommandHandler(t *testing.T) {
	t.Run("should package successfully", func(t *testing.T) {
		type arg struct {
			source         *string
			target         *string
			expectedSource string
			expectedTarget string
		}

		folder := "/x"
		arguments := []arg{
			{
				source:         &folder,
				target:         &folder,
				expectedSource: folder,
				expectedTarget: folder,
			},
			{
				source:         nil,
				target:         nil,
				expectedSource: ".",
				expectedTarget: "..",
			},
		}

		for _, argument := range arguments {
			controller := gomock.NewController(t)
			mockCommander := operator.NewMockOperator(controller)

			expectedCall := &operator.CompressDirectoryOptions{
				SourceDirectory: &argument.expectedSource,
				TargetDirectory: &argument.expectedTarget,
				Compressor:      compressor.GetInstance(),
				Logger:          logger.NewLogger(),
			}

			mockCommander.EXPECT().CompressDirectory(gomock.Eq(expectedCall))

			command := CreatePackageCommand(mockCommander)

			if argument.source != nil || argument.target != nil {
				command.SetArgs(getPackageTestArguments(*argument.source, *argument.target))
			}

			err := command.Execute()
			require.Nil(t, err)
		}
	})
}

func getPackageTestArguments(source, target string) []string {
	return []string{"-s", source, "-t", target}
}
