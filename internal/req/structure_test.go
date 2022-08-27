// +build unit

package req

import (
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
	projectStructure1 = model.ProjectStructureData{
		Name:      "Project -1",
		Reference: "go.dev",
		URL:       "https://project1.com",
	}
	projectStructure2 = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
	}
	testProjectData = []*model.ProjectStructureData{&projectStructure1, &projectStructure2}
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

		options := make([]prompter.Option, 0)
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

		input, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, input)
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

		input, err := p.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, prompter.EmptyList)

		require.Nil(t, input)
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
