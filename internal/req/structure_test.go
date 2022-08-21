package req

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	projectStructure1 = lister.ProjectStructureData{
		Name:      "Project -1",
		Reference: "go.dev",
		URL:       "https://project1.com",
	}
	projectStructure2 = lister.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
	}
	testProjectData = []*lister.ProjectStructureData{&projectStructure1, &projectStructure2}

	testDataWithOneStructure = []*lister.ProjectStructureData{&projectStructure1}
)

func TestStructure_AskForInput(t *testing.T) {
	t.Run("should ask for selection for project", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompts.NewMockPrompter(controller)
		mockUncompressor := uncompressor.NewMockUncompressor(controller)

		expectedValue := []*prompts.ListOption{
			{
				DisplayText: testProjectData[0].String(),
				ReturnVal:   testProjectData[0],
			},
			{
				DisplayText: testProjectData[1].String(),
				ReturnVal:   testProjectData[1],
			},
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Eq(SelectProjectTypeDirection), gomock.Eq(expectedValue)).
			Return(testProjectData[0]).
			Times(1)

		p := &ProjectStructureRequirement{
			ProjectsData: testProjectData,
			P:            mockPrompter,
			U:            mockUncompressor,
		}

		input, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, input)

		task := input.(*projectStructureTask)

		require.EqualValues(t, task.ProjectStructure, testProjectData[0])
		require.NotNil(t, task.U)
	})

	t.Run("should return error if the project data structure list is empty", func(t *testing.T) {
		p := &ProjectStructureRequirement{
			ProjectsData: nil,
		}

		input, err := p.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrProductStructureListIsEmpty)

		require.Nil(t, input)
	})

	t.Run("should return first item if there is only one project structure", func(t *testing.T) {
		p := &ProjectStructureRequirement{
			ProjectsData: testDataWithOneStructure,
		}

		input, err := p.AskForInput()

		require.Nil(t, err)
		require.NotNil(t, input)

		task := input.(*projectStructureTask)
		require.EqualValues(t, task.ProjectStructure, testDataWithOneStructure[0])
	})
}

func TestStructure_Complete(t *testing.T) {
	//TODO write test
}
