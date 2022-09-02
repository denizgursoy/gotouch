// +build unit

package req

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	questions = []*model.Question{
		{
			Direction:         "question 1",
			CanSkip:           false,
			CanSelectMultiple: false,
			Options:           nil,
		},
		{
			Direction:         "question 2",
			CanSkip:           false,
			CanSelectMultiple: false,
			Options:           nil,
		},
	}
	projectStructure1 = model.ProjectStructureData{
		Name:      "Project -1",
		Reference: "go.dev",
		URL:       "https://project1.com",
		Questions: questions,
	}
	projectStructure2 = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
	}

	projectStructureWithValues = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
		Values: map[interface{}]interface{}{
			"x": "",
		},
	}

	testProjectData                = []*model.ProjectStructureData{&projectStructure1, &projectStructure2}
	testProjectDataWithValuesField = []*model.ProjectStructureData{&projectStructureWithValues, &projectStructure2}
)

func TestStructure_AskForInput(t *testing.T) {
	t.Run("should ask for selection for project", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockUncompressor := compressor.NewMockCompressor(controller)
		mockManager := manager.NewMockManager(controller)
		mockExecutor := executor.NewMockExecutor(controller)
		mockStore := store.GetInstance()

		options := make([]fmt.Stringer, 0)
		for _, datum := range testProjectData {
			options = append(options, datum)
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(SelectProjectTypeDirection), gomock.Eq(options)).
			Return(testProjectData[0], nil).
			Times(1)

		p := &ProjectStructureRequirement{
			ProjectsData: testProjectData,
			Prompter:     mockPrompter,
			Compressor:   mockUncompressor,
			Logger:       logger.NewLogger(),
			Executor:     mockExecutor,
			Manager:      mockManager,
			Store:        mockStore,
		}

		tasks, requirements, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, tasks)

		require.Len(t, tasks, 1)
		require.IsType(t, tasks[0], &projectStructureTask{})
		require.IsType(t, tasks[0].(*projectStructureTask).ProjectStructure, testProjectData[0])

		actualQuestions := make([]*model.Question, 0)
		for _, requirement := range requirements {
			actualQuestions = append(actualQuestions, &requirement.(*QuestionRequirement).Question)
		}
		require.Equal(t, questions, actualQuestions)
	})

	t.Run("should add template requirement if value is not nil", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockUncompressor := compressor.NewMockCompressor(controller)
		mockManager := manager.NewMockManager(controller)
		mockExecutor := executor.NewMockExecutor(controller)
		mockStore := store.GetInstance()

		options := make([]fmt.Stringer, 0)
		for _, datum := range testProjectData {
			options = append(options, datum)
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(SelectProjectTypeDirection), gomock.Eq(options)).
			Return(testProjectDataWithValuesField[0], nil).
			Times(1)

		p := &ProjectStructureRequirement{
			ProjectsData: testProjectData,
			Prompter:     mockPrompter,
			Compressor:   mockUncompressor,
			Logger:       logger.NewLogger(),
			Executor:     mockExecutor,
			Manager:      mockManager,
			Store:        mockStore,
		}

		_, requirements, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, requirements)

		require.Len(t, requirements, 1)
		require.IsType(t, requirements[0], &templateRequirement{})

	})

	t.Run("should return error from the prompt", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockCompressor := compressor.NewMockCompressor(controller)
		mockManager := manager.NewMockManager(controller)
		mockExecutor := executor.NewMockExecutor(controller)
		mockLogger := logger.NewLogger()
		mockStore := store.GetInstance()

		p := &ProjectStructureRequirement{
			Prompter:   mockPrompter,
			Compressor: mockCompressor,
			Manager:    mockManager,
			Logger:     mockLogger,
			Executor:   mockExecutor,
			Store:      mockStore,
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Any(), gomock.Any()).
			Return(nil, prompter.EmptyList).
			Times(1)

		tasks, _, err := p.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, prompter.EmptyList)

		require.Empty(t, tasks)
	})

}

func TestStructure_Complete(t *testing.T) {
	t.Run("should call uncompress with the URL", func(t *testing.T) {
		type testCase struct {
			ProjectName   string
			DirectoryName string
		}

		controller := gomock.NewController(t)

		defer controller.Finish()

		testCases := []testCase{
			{ProjectName: testProjectName, DirectoryName: testProjectName},
			{ProjectName: testUrlName, DirectoryName: testProjectName},
		}
		for _, testCase := range testCases {
			mockUncompressor := compressor.NewMockCompressor(controller)
			mockManager := manager.NewMockManager(controller)
			mockLogger := logger.NewLogger()
			mockExecutor := executor.NewMockExecutor(controller)
			mockStore := store.NewMockStore(controller)

			mockStore.EXPECT().GetValue(store.ModuleName).Return(testCase.ProjectName).Times(1)

			mockUncompressor.
				EXPECT().
				UncompressFromUrl(gomock.Eq(projectStructure1.URL))

			mockManager.
				EXPECT().
				EditGoModule().
				AnyTimes()

			p := &projectStructureTask{
				ProjectStructure: &projectStructure1,
				Compressor:       mockUncompressor,
				Manager:          mockManager,
				Logger:           mockLogger,
				Executor:         mockExecutor,
				Store:            mockStore,
			}
			err := p.Complete()
			require.Nil(t, err)
		}

	})
}
