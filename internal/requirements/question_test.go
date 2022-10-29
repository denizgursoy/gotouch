package requirements

import (
	"errors"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/langs"
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
	choice = model.Choice{
		Choice:       "112322",
		Dependencies: []interface{}{dependency1, dependency2},
		Files:        []*model.File{&file1, &file2},
		Values: map[string]interface{}{
			"X": "sds",
		},
	}

	yesNoQuestion = model.Question{
		Direction:         "yes no question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Choices:           []*model.Choice{&choice},
	}

	multipleChoiceQuestion = model.Question{
		Direction:         "yes no question",
		CanSkip:           false,
		CanSelectMultiple: false,
		Choices:           []*model.Choice{&choice, &choice},
	}

	multipleChoiceQuestionWithSkip = model.Question{
		Direction:         "yes no question",
		CanSkip:           true,
		CanSelectMultiple: false,
		Choices:           []*model.Choice{&choice, &choice},
	}
	promptErr = errors.New("prompt-err")
)

func TestQuestionRequirement_AskForInput(t *testing.T) {
	t.Run("should call yes/no question and return 4 tasks if canskip is true and there is only one choice", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, yesNoQuestion)
		defer controller.Finish()

		requirement.Store.(*store.MockStore).EXPECT().StoreValues(gomock.Eq(yesNoQuestion.Choices[0].Values))
		requirement.Prompter.(*prompter.MockPrompter).EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(true, nil).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Nil(t, requirements)

		require.Len(t, task, len(yesNoQuestion.Choices[0].Dependencies)+len(yesNoQuestion.Choices[0].Files))
	})

	t.Run("should return no task if no is selected", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, yesNoQuestion)
		defer controller.Finish()

		requirement.Prompter.(*prompter.MockPrompter).EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(false, nil).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Nil(t, requirements)

		require.Len(t, task, 0)
	})

	t.Run("should return error if prompt returns error", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, yesNoQuestion)
		defer controller.Finish()

		requirement.Prompter.(*prompter.MockPrompter).EXPECT().AskForYesOrNo(gomock.Eq(yesNoQuestion.Direction)).Return(false, promptErr).Times(1)

		task, requirements, err := requirement.AskForInput()
		require.ErrorIs(t, promptErr, err)
		require.Nil(t, requirements)
		require.Nil(t, task)
	})

	t.Run("should select from list if there is more than 1 choice", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, multipleChoiceQuestion)
		defer controller.Finish()

		choices := make([]fmt.Stringer, 0)
		for _, choice := range requirement.Question.Choices {
			choices = append(choices, choice)
		}

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleChoiceQuestion.Direction), gomock.Eq(choices)).
			Return(multipleChoiceQuestion.Choices[0], nil).
			Times(1)

		requirement.Store.(*store.MockStore).EXPECT().StoreValues(gomock.Eq(multipleChoiceQuestion.Choices[0].Values))
		_, _, _ = requirement.AskForInput()
	})

	t.Run("should add none of above choice if canskip is true", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, multipleChoiceQuestionWithSkip)
		defer controller.Finish()

		choices := make([]fmt.Stringer, 0)
		for _, choice := range requirement.Question.Choices {
			choices = append(choices, choice)
		}
		choices = append(choices, noneOfAboveChoice)

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleChoiceQuestionWithSkip.Direction), gomock.Eq(choices)).
			Return(multipleChoiceQuestionWithSkip.Choices[0], nil).
			Times(1)

		requirement.Store.(*store.MockStore).EXPECT().StoreValues(gomock.Eq(multipleChoiceQuestion.Choices[0].Values))
		_, _, _ = requirement.AskForInput()
	})

	t.Run("should return error if select from list returns errors", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, multipleChoiceQuestionWithSkip)
		defer controller.Finish()

		choices := make([]fmt.Stringer, 0)
		for _, choice := range requirement.Question.Choices {
			choices = append(choices, choice)
		}
		choices = append(choices, noneOfAboveChoice)

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Eq(multipleChoiceQuestionWithSkip.Direction), gomock.Eq(choices)).
			Return(nil, promptErr).
			Times(1)

		_, _, err := requirement.AskForInput()
		require.NotNil(t, err)
		require.ErrorIs(t, err, promptErr)
	})
}

func getTestQuestionRequirement(t *testing.T, question model.Question) (*QuestionRequirement, *gomock.Controller) {
	controller := gomock.NewController(t)
	mockExecutor := executor.NewMockExecutor(controller)
	mockPrompter := prompter.NewMockPrompter(controller)
	mockManager := manager.NewMockManager(controller)
	mockStore := store.NewMockStore(controller)
	mockChecker := langs.NewMockChecker(controller)

	return &QuestionRequirement{
		Question:        question,
		Prompter:        mockPrompter,
		Logger:          logger.NewLogger(),
		Executor:        mockExecutor,
		Manager:         mockManager,
		Store:           mockStore,
		LanguageChecker: mockChecker,
	}, controller
}
