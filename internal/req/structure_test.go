package req

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
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

		mockPrompter := prompts.NewMockPrompter(controller)
		mockUncompressor := uncompressor.NewMockUncompressor(controller)

		options := make([]prompts.Option, 0)
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
			Uncompressor: mockUncompressor,
		}

		input, err := p.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, input)

		task := input.(*projectStructureTask)

		require.EqualValues(t, task.ProjectStructure, testProjectData[0])
		require.NotNil(t, task.U)
	})

	t.Run("should return error from the prompt", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompts.NewMockPrompter(controller)

		p := &ProjectStructureRequirement{
			Prompter: mockPrompter,
		}

		mockPrompter.
			EXPECT().
			AskForSelectionFromList(gomock.Any(), gomock.Any()).
			Return(nil, prompts.ErrProductStructureListIsEmpty).
			Times(1)

		input, err := p.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, prompts.ErrProductStructureListIsEmpty)

		require.Nil(t, input)
	})

}

func TestStructure_Complete(t *testing.T) {
	//TODO write test
}
