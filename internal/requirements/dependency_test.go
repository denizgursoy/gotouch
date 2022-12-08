package requirements

import (
	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
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
