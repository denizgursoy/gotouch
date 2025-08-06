package requirements

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/denizgursoy/gotouch/internal/langs"
)

func Test_cleanupTask_Complete(t *testing.T) {
	t.Run("should call setup of checker", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockChecker := langs.NewMockChecker(controller)

		testCleanupTask := &cleanupTask{
			LanguageChecker: mockChecker,
		}

		mockChecker.EXPECT().CleanUp().Times(1)

		err := testCleanupTask.Complete(context.Background())

		require.Nil(t, err)
	})
}
