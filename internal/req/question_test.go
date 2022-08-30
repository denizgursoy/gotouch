package req

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	dependency1 = "2132"
	dependency2 = "2132"
	file1       = model.File{
		Url:     "",
		Content: "",
		Path:    "",
	}
	file2 = model.File{
		Url:     "",
		Content: "",
		Path:    "",
	}
	yesNoQuestion = model.Question{
		Direction:         "yes no question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options: []*model.Option{
			{
				Answer:       "112322",
				Dependencies: []*string{&dependency1, &dependency2},
				Files:        []*model.File{&file1, &file2},
			},
		},
	}
)

func TestQuestionRequirement_AskForInput(t *testing.T) {
	t.Run("should call yes/no question and return 4 tasks if canskip is true and there is only one choice", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
		}

		mockPrompter.EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(true, nil).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Nil(t, requirements)

		require.Len(t, task, len(yesNoQuestion.Options[0].Dependencies)+len(yesNoQuestion.Options[0].Files))

	})

	t.Run("should return no task if no is selected", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
		}

		mockPrompter.EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(false, nil).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Nil(t, requirements)

		require.Len(t, task, 0)
	})

	t.Run("should return error if prompt returns error", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
		}

		promptErr := errors.New("prompt-err")
		mockPrompter.EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(false, promptErr).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.ErrorIs(t, promptErr, err)
		require.Nil(t, requirements)
		require.Nil(t, task)
	})
}
