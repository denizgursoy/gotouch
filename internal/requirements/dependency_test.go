package requirements

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestComplete(t *testing.T) {
	t.Run("should call lang checker", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		checker := langs.NewMockChecker(controller)
		dependency := "test-dependency"

		checker.EXPECT().GetDependency(gomock.Eq(dependency))

		task := dependencyTask{
			Dependency:      dependency,
			LanguageChecker: checker,
		}
		err := task.Complete()
		require.Nil(t, err)
	})
}
