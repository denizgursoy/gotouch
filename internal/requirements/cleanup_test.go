package requirements

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/langs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_cleanupTask_Complete(t *testing.T) {
	t.Run("should call setup of checker", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockChecker := langs.NewMockChecker(controller)

		testCleanupTask := &cleanupTask{
			LanguageChecker: mockChecker,
		}

		mockChecker.EXPECT().CleanUp().Times(1)

		err := testCleanupTask.Complete()

		require.Nil(t, err)
	})
}
