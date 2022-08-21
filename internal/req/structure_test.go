// +build unit

package req

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
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
		}

		input, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, input)

		task := input.(*projectStructureTask)

		require.EqualValues(t, task.ProjectStructure, testProjectData[0])
		require.NotNil(t, task.Uncompressor)
	})

	t.Run("should return error from the prompt", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)

		p := &ProjectStructureRequirement{
			Prompter: mockPrompter,
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Any(), gomock.Any()).
			Return(nil, prompter.ErrProductStructureListIsEmpty).
			Times(1)

		input, err := p.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, prompter.ErrProductStructureListIsEmpty)

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

			mockUncompressor.
				EXPECT().
				UncompressFromUrl(gomock.Eq(projectStructure1.URL), gomock.Eq(testCase.DirectoryName))

			mockManager.
				EXPECT().
				EditGoModule(gomock.Eq(testCase.ProjectName), gomock.Eq(testCase.DirectoryName))

			p := &projectStructureTask{
				ProjectStructure: &projectStructure1,
				Uncompressor:     mockUncompressor,
				Manager:          mockManager,
			}
			actualData, err := p.Complete(testCase.ProjectName)
			require.Nil(t, err)
			require.Nil(t, actualData)
		}

	})
}
