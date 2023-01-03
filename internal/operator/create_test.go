package operator

import (
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/requirements"
	"github.com/denizgursoy/gotouch/internal/store"
)

var project1 = model.ProjectStructureData{
	Name:      "sds",
	Reference: "sds",
	URL:       "sds",
}

func TestCreateNewProject(t *testing.T) {
	t.Run("should call Operator with all requirements", func(t *testing.T) {
		options, controller := createTestNewProjectOptions(t, nil)
		defer controller.Finish()

		expectedProjectData := []*model.ProjectStructureData{&project1}
		options.Lister.(*lister.MockLister).
			EXPECT().
			GetProjectList(nil).
			Return(expectedProjectData, nil).
			Times(1)

		options.Executor.(*executor.MockExecutor).
			EXPECT().
			Execute(gomock.Any()).
			Do(func(arg any) {
				execRequirements := arg.(executor.Requirements)
				require.Len(t, execRequirements, 1)
				structure := execRequirements[0].(*requirements.ProjectStructureRequirement)

				require.NotNil(t, structure.Compressor)
				require.NotNil(t, structure.Manager)
				require.NotNil(t, structure.Prompter)

				require.IsType(t, (*requirements.ProjectStructureRequirement)(nil), structure)
				require.EqualValues(t, expectedProjectData, structure.ProjectsData)
			})

		err := GetInstance().CreateNewProject(&options)

		require.Nil(t, err)
	})
}

func Test_isValid(t *testing.T) {
	options, controller := createTestNewProjectOptions(t, nil)
	defer controller.Finish()

	got := isValid(&options)
	require.Nil(t, got)
}

func Test_isValid_PathTest(t *testing.T) {
	options, controller := createTestNewProjectOptions(t, nil)
	defer controller.Finish()

	t.Run("should return no error if path is nil", func(t *testing.T) {
		arg := options
		err := isValid(&arg)
		require.Nil(t, err)
	})

	t.Run("should return error if path does not end with yaml", func(t *testing.T) {
		arg := options
		path := "test.zaml"
		arg.Path = &path
		err := isValid(&arg)
		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrNotYamlFile)
	})

	t.Run("should return no error if yaml file exists", func(t *testing.T) {
		arg := options
		path := "../testdata/input.yaml"
		arg.Path = &path
		err := isValid(&arg)
		require.Nil(t, err)
	})

	t.Run("should return error if file does not exists", func(t *testing.T) {
		arg := options
		path := "./testdata/input2.yaml"
		arg.Path = &path
		err := isValid(&arg)
		require.NotNil(t, err)
		require.ErrorIs(t, err, ErrNotValidUrlOrFilePath)
	})
}

func createTestNewProjectOptions(t *testing.T, path *string) (CreateNewProjectOptions, *gomock.Controller) {
	controller := gomock.NewController(t)

	mockLister := lister.NewMockLister(controller)
	mockPrompter := prompter.NewMockPrompter(controller)
	newMockManager := manager.NewMockManager(controller)
	mockCompressor := compressor.NewMockCompressor(controller)
	mockExecutor := executor.NewMockExecutor(controller)
	mockLogger := logger.NewLogger()
	mockStore := store.GetInstance()
	mockCloner := cloner.NewMockCloner(controller)
	mockRunner := commandrunner.NewMockRunner(controller)

	return CreateNewProjectOptions{
		Lister:        mockLister,
		Prompter:      mockPrompter,
		Manager:       newMockManager,
		Compressor:    mockCompressor,
		Executor:      mockExecutor,
		Logger:        mockLogger,
		Store:         mockStore,
		Path:          nil,
		Cloner:        mockCloner,
		CommandRunner: mockRunner,
	}, controller

}
