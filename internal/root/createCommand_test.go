// +build unit

package root

import (
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/req"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	project1 = model.ProjectStructureData{
		Name:      "sds",
		Reference: "sds",
		URL:       "sds",
	}
)

func TestCreateNewProject(t *testing.T) {
	t.Run("should call executor with all requirements", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockLister := lister.NewMockLister(controller)
		mockExecutor := executor.NewMockExecutor(controller)

		expectedProjectData := []*model.ProjectStructureData{&project1}
		mockLister.
			EXPECT().
			GetProjectList(nil).
			Return(expectedProjectData, nil).
			Times(1)

		mockExecutor.
			EXPECT().
			Execute(gomock.Any()).
			Do(func(arg interface{}) {
				requirements := arg.(executor.Requirements)
				require.Len(t, requirements, 2)
				name := requirements[0].(*req.ProjectNameRequirement)
				structure := requirements[1].(*req.ProjectStructureRequirement)

				require.IsType(t, (*req.ProjectNameRequirement)(nil), name)
				require.IsType(t, (*req.ProjectStructureRequirement)(nil), structure)
				require.EqualValues(t, expectedProjectData, structure.ProjectsData)
			})

		opts := &CreateCommandOptions{
			lister:   mockLister,
			executor: mockExecutor,
		}
		err := CreateNewProject(opts)

		require.Nil(t, err)
	})
}
