package langs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	dependencyWithoutVersion = "github.com/labstack/echo/v4"
	latestDependencies       = dependencyWithoutVersion + atSign + latestVersion
	dependencyWithVersion    = "github.com/labstack/echo/v4@v1.2.3"
)

//func Test_GetDependency(t *testing.T) {
//	t.Run("should call add dependency method with latest if version is not present", func(t *testing.T) {
//		controller := gomock.NewController(t)
//
//		type arg struct {
//			Dependency string
//			CallValue  CommandData
//		}
//
//		dependencies := []arg{
//			{
//				Dependency: dependencyWithoutVersion,
//				CallValue: CommandData{
//					Command: "go",
//					Args:    []string{"get", latestDependencies},
//				},
//			},
//			{
//				Dependency: dependencyWithVersion,
//				CallValue: CommandData{
//					Command: "go",
//					Args:    []string{"get", dependencyWithVersion},
//				},
//			},
//		}
//
//		for _, testArg := range dependencies {
//			mockExecutor := executor.NewMockExecutor(controller)
//
//			mockExecutor.
//				EXPECT().
//				RunCommand(gomock.Eq(&testArg.CallValue)).
//				Return(nil).
//				Times(1)
//
//			checker := NewGolangSetupChecker(logger.NewLogger(), store.NewMockStore(controller))
//			err := checker.GetDependency(testArg.Dependency)
//			require.Nil(t, err)
//		}
//	})
//
//	t.Run("should return error if dependency is not installed ", func(t *testing.T) {
//		controller := gomock.NewController(t)
//		mockExecutor := executor.NewMockExecutor(controller)
//
//		commandError := errors.New("commander error")
//		mockExecutor.
//			EXPECT().
//			RunCommand(gomock.Any()).
//			Return(commandError).
//			Times(1)
//
//		checker := NewGolangSetupChecker(logger.NewLogger(), store.NewMockStore(controller))
//		err := checker.GetDependency(dependencyWithoutVersion)
//
//		require.ErrorIs(t, commandError, err)
//	})
//
//	t.Run("should return error if string is empty", func(t *testing.T) {
//		controller := gomock.NewController(t)
//		mockExecutor := executor.NewMockExecutor(controller)
//
//		mockExecutor.EXPECT().RunCommand(gomock.Any()).AnyTimes()
//
//		checker := NewGolangSetupChecker(logger.NewLogger(), store.NewMockStore(controller))
//		err := checker.GetDependency("   ")
//
//		require.ErrorIs(t, ErrDependencyIsNotValid, err)
//	})
//}

func TestDisplayName(t *testing.T) {
	t.Run("should include reference if exists", func(t *testing.T) {
		version := "1234"
		url := "*"
		dependency := Dependency{Url: &url, Version: &version}
		require.EqualValues(t, url+"@"+version, dependency.String())
	})
}
