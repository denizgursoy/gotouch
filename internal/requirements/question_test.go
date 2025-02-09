package requirements

import (
	"errors"
	"fmt"
	"testing"

	"github.com/denizgursoy/gotouch/internal/langs"

	"github.com/denizgursoy/gotouch/internal/store"

	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	dependency1 = "dependency-1"
	dependency2 = "dependency-2"
	dependency3 = "dependency-3"
	dependency4 = "dependency-4"
	dependency5 = "dependency-5"
	file1       = model.File{
		Url:          "file url1",
		Content:      "content-1",
		PathFromRoot: "path-1",
	}
	file2 = model.File{
		Url:          "file url2",
		Content:      "content-2",
		PathFromRoot: "path2",
	}
	choice = model.Choice{
		Choice: "choice 1",
		Resources: model.Resources{
			Dependencies: []any{dependency1, dependency2},
			Files:        []*model.File{&file1, &file2},
			Values: map[string]any{
				"X": "sds",
			},
			CustomValues: map[string]any{
				"foo": "bar",
			},
		},
	}
	choice2 = model.Choice{
		Choice: "choice 2",
		Resources: model.Resources{
			Dependencies: []any{dependency3},
			Files:        []*model.File{&file1, &file2},
			Values: map[string]any{
				"Y": "sds",
			},
			CustomValues: map[string]any{
				"foo2": "bar2",
			},
		},
	}
	choice3 = model.Choice{
		Choice: "choice 3",
		Resources: model.Resources{
			Dependencies: []any{dependency4, dependency5},
			Files:        []*model.File{&file2},
			Values: map[string]any{
				"Z": "sds",
			},
			CustomValues: map[string]any{
				"foo3": "bar3",
			},
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
	multipleChoiceQuestionWithMultiSelection = model.Question{
		Direction:         "yes no question",
		CanSelectMultiple: true,
		Choices:           []*model.Choice{&choice, &choice2, &choice3},
	}
	promptErr = errors.New("prompt-err")
)

func TestQuestionRequirement_AskForInput(t *testing.T) {
	t.Run("should call yes/no question and return 4 tasks if canskip is true and there is only one choice", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, yesNoQuestion)
		defer controller.Finish()

		requirement.Store.(*store.MockStore).EXPECT().AddValues(gomock.Eq(yesNoQuestion.Choices[0].Values))
		requirement.Store.(*store.MockStore).EXPECT().AddCustomValues(gomock.Eq(yesNoQuestion.Choices[0].CustomValues))
		requirement.Store.(*store.MockStore).EXPECT().AddDependency(gomock.Eq(dependency1))
		requirement.Store.(*store.MockStore).EXPECT().AddDependency(gomock.Eq(dependency2))
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

		requirement.Store.(*store.MockStore).EXPECT().AddValues(gomock.Eq(multipleChoiceQuestion.Choices[0].Values))
		requirement.Store.(*store.MockStore).EXPECT().AddCustomValues(gomock.Eq(multipleChoiceQuestion.Choices[0].CustomValues))
		requirement.Store.(*store.MockStore).EXPECT().AddDependency(gomock.Any()).AnyTimes()

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

		requirement.Store.(*store.MockStore).EXPECT().AddValues(gomock.Eq(multipleChoiceQuestion.Choices[0].Values))
		requirement.Store.(*store.MockStore).EXPECT().AddCustomValues(gomock.Eq(multipleChoiceQuestion.Choices[0].CustomValues))
		requirement.Store.(*store.MockStore).EXPECT().AddDependency(gomock.Any()).AnyTimes()

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

	t.Run("should ask for multiple choice from prompter", func(t *testing.T) {
		requirement, controller := getTestQuestionRequirement(t, multipleChoiceQuestionWithMultiSelection)
		defer controller.Finish()

		choices := make([]fmt.Stringer, 0)
		for _, choice := range requirement.Question.Choices {
			choices = append(choices, choice)
		}

		selectedChoices := make([]any, 0)
		selectedChoices = append(selectedChoices, choices[0], choices[2])

		for _, selectedChoice := range selectedChoices {
			chc := selectedChoice.(*model.Choice)

			requirement.Store.(*store.MockStore).EXPECT().AddValues(gomock.Eq(chc.Values)).Times(1)
			requirement.Store.(*store.MockStore).EXPECT().AddCustomValues(gomock.Eq(chc.CustomValues)).Times(1)

			for _, dependency := range chc.Dependencies {
				requirement.Store.(*store.MockStore).
					EXPECT().
					AddDependency(gomock.Eq(dependency)).
					AnyTimes()
			}
		}

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForMultipleSelectionFromList(gomock.Eq(multipleChoiceQuestionWithSkip.Direction), gomock.Eq(choices)).
			Return(selectedChoices, nil).
			Times(1)

		tasks, _, err := requirement.AskForInput()
		require.Nil(t, err)
		require.Len(t, tasks, 7)

		require.Equal(t, choices[0].(*model.Choice).Dependencies[0], tasks[0].(*dependencyTask).Dependency)
		require.Equal(t, choices[0].(*model.Choice).Dependencies[1], tasks[1].(*dependencyTask).Dependency)
		require.Equal(t, *choices[0].(*model.Choice).Files[0], tasks[2].(*fileTask).File)
		require.Equal(t, *choices[0].(*model.Choice).Files[1], tasks[3].(*fileTask).File)

		require.Equal(t, choices[2].(*model.Choice).Dependencies[0], tasks[4].(*dependencyTask).Dependency)
		require.Equal(t, choices[2].(*model.Choice).Dependencies[1], tasks[5].(*dependencyTask).Dependency)
		require.Equal(t, *choices[2].(*model.Choice).Files[0], tasks[6].(*fileTask).File)
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
