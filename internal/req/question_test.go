package req

import (
	"errors"
	"fmt"
	"testing"

	"github.com/denizgursoy/gotouch/internal/store"

	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	dependency1 = "2132"
	dependency2 = "2132"
	file1       = model.File{
		Url:          "",
		Content:      "",
		PathFromRoot: "",
	}
	file2 = model.File{
		Url:          "",
		Content:      "",
		PathFromRoot: "",
	}
	option = model.Option{
		Choice:       "112322",
		Dependencies: []*string{&dependency1, &dependency2},
		Files:        []*model.File{&file1, &file2},
		Values: map[string]interface{}{
			"X": "sds",
		},
	}

	yesNoQuestion = model.Question{
		Direction:         "yes no question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options:           []*model.Option{&option},
	}

	multipleOptionQuestion = model.Question{
		Direction:         "yes no question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Options:           []*model.Option{&option, &option},
	}

	multipleOptionQuestionWithSkip = model.Question{
		Direction:         "yes no question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Options:           []*model.Option{&option, &option},
	}
	promptErr = errors.New("prompt-err")
)

func TestQuestionRequirement_AskForInput(t *testing.T) {
	t.Run("should call yes/no question and return 4 tasks if canskip is true and there is only one choice", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
		}

		mockStore.EXPECT().StoreValues(gomock.Eq(yesNoQuestion.Options[0].Values))
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
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
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
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: yesNoQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
		}

		mockPrompter.EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(false, promptErr).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.ErrorIs(t, promptErr, err)
		require.Nil(t, requirements)
		require.Nil(t, task)
	})

	t.Run("should select from list if there is more than 1 option", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: multipleOptionQuestion,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
		}

		options := make([]fmt.Stringer, 0)
		for _, option := range requirement.Question.Options {
			options = append(options, option)
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleOptionQuestion.Direction), gomock.Eq(options)).
			Return(multipleOptionQuestion.Options[0], nil).
			Times(1)

		mockStore.EXPECT().StoreValues(gomock.Eq(multipleOptionQuestion.Options[0].Values))
		_, _, _ = requirement.AskForInput()
	})

	t.Run("should add none of above option if canskip is true", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: multipleOptionQuestionWithSkip,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
		}

		options := make([]fmt.Stringer, 0)
		for _, option := range requirement.Question.Options {
			options = append(options, option)
		}
		options = append(options, noneOfAboveOption)

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleOptionQuestionWithSkip.Direction), gomock.Eq(options)).
			Return(multipleOptionQuestionWithSkip.Options[0], nil).
			Times(1)

		mockStore.EXPECT().StoreValues(gomock.Eq(multipleOptionQuestion.Options[0].Values))
		_, _, _ = requirement.AskForInput()
	})

	t.Run("should return error if select from list returns errors", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockExecutor := executor.NewMockExecutor(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockStore := store.NewMockStore(controller)

		requirement := QuestionRequirement{
			Question: multipleOptionQuestionWithSkip,
			Prompter: mockPrompter,
			Logger:   logger.NewLogger(),
			Executor: mockExecutor,
			Manager:  mockManager,
			Store:    mockStore,
		}

		options := make([]fmt.Stringer, 0)
		for _, option := range requirement.Question.Options {
			options = append(options, option)
		}
		options = append(options, noneOfAboveOption)

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleOptionQuestionWithSkip.Direction), gomock.Eq(options)).
			Return(nil, promptErr).
			Times(1)

		_, _, err := requirement.AskForInput()
		require.NotNil(t, err)
		require.ErrorIs(t, err, promptErr)
	})
}
