package requirements

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
)

func TestAskForInput(t *testing.T) {
	t.Run("should not prompt to user if values are empty", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockStore := store.NewMockStore(controller)

		mockStore.EXPECT().GetCustomValues().Times(1)

		requirement := templateRequirement{
			Prompter: mockPrompter,
			Store:    mockStore,
		}

		input, requirements, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Nil(t, requirements)
		require.Len(t, input, 1)
		require.IsType(t, input[0], &templateTask{})
	})
}
