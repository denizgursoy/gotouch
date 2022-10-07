//go:build unit
// +build unit

package req

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	dependencyWithoutVersion = "github.com/labstack/echo/v4"
	latestDependencies       = dependencyWithoutVersion + atSign + latestVersion
	dependencyWithVersion    = "github.com/labstack/echo/v4@v1.2.3"
)

func Test_dependencyTask_Complete(t *testing.T) {
	t.Run("should call add dependency method with latest if version is not present", func(t *testing.T) {
		controller := gomock.NewController(t)

		type arg struct {
			Dependency string
			CallValue  executor.CommandData
		}

		dependencies := []arg{
			{
				Dependency: dependencyWithoutVersion,
				CallValue: executor.CommandData{
					Command: "go",
					Args:    []string{"get", latestDependencies},
				},
			},
			{
				Dependency: dependencyWithVersion,
				CallValue: executor.CommandData{
					Command: "go",
					Args:    []string{"get", dependencyWithVersion},
				},
			},
		}

		for _, testArg := range dependencies {
			mockExecutor := executor.NewMockExecutor(controller)

			mockExecutor.
				EXPECT().
				RunCommand(gomock.Eq(&testArg.CallValue)).
				Return(nil).
				Times(1)

			task := dependencyTask{
				Dependency: testArg.Dependency,
				Logger:     logger.NewLogger(),
				Executor:   mockExecutor,
			}

			err := task.Complete()
			require.Nil(t, err)
		}

	})

	t.Run("should return error if dependency is not installed ", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)

		commandError := errors.New("commander error")
		mockExecutor.
			EXPECT().
			RunCommand(gomock.Any()).
			Return(commandError).
			Times(1)

		task := dependencyTask{
			Dependency: dependencyWithoutVersion,
			Logger:     logger.NewLogger(),
			Executor:   mockExecutor,
		}

		err := task.Complete()
		require.ErrorIs(t, commandError, err)
	})

	t.Run("should return error if string is empty", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)

		mockExecutor.EXPECT().RunCommand(gomock.Any()).AnyTimes()

		task := dependencyTask{
			Dependency: "   ",
			Logger:     logger.NewLogger(),
			Executor:   mockExecutor,
		}

		err := task.Complete()
		require.ErrorIs(t, ErrDependencyIsNotValid, err)
	})
}
