package req

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	dependencyWithoutVersion = "github.com/labstack/echo/v4"
	dependencyWithVersion    = "github.com/labstack/echo/v4@latest"
)

func Test_dependencyTask_Complete(t *testing.T) {
	t.Run("should call add dependency method", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockManager := manager.NewMockManager(controller)

		mockManager.
			EXPECT().
			AddDependency(dependencyWithoutVersion).
			Return(nil).
			Times(1)

		task := dependencyTask{
			Manager:    mockManager,
			Dependency: dependencyWithoutVersion,
		}

		complete, err := task.Complete(nil)
		require.Nil(t, err)
		require.Nil(t, complete)
	})

	t.Run("should return error if dependency is not installed ", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockManager := manager.NewMockManager(controller)

		dependencyError := errors.New("dependency error")
		mockManager.
			EXPECT().
			AddDependency(dependencyWithoutVersion).
			Return(dependencyError).
			Times(1)

		task := dependencyTask{
			Manager:    mockManager,
			Dependency: dependencyWithoutVersion,
		}

		complete, err := task.Complete(nil)
		require.ErrorIs(t, dependencyError, err)
		require.Nil(t, complete)
	})
}
