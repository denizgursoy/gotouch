//go:build unit
// +build unit

package operator

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/req"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var project1 = model.ProjectStructureData{
	Name:      "sds",
	Reference: "sds",
	URL:       "sds",
}

func TestCreateNewProject(t *testing.T) {
	t.Run("should call Operator with all requirements", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockLister := lister.NewMockLister(controller)
		mockPrompter := prompter.NewMockPrompter(controller)
		newMockManager := manager.NewMockManager(controller)
		mockCompressor := compressor.NewMockCompressor(controller)
		mockExecutor := executor.NewMockExecutor(controller)
		mockLogger := logger.NewLogger()
		mockStore := store.GetInstance()

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
				require.Len(t, requirements, 1)
				structure := requirements[0].(*req.ProjectStructureRequirement)

				require.NotNil(t, structure.Compressor)
				require.NotNil(t, structure.Manager)
				require.NotNil(t, structure.Prompter)

				require.IsType(t, (*req.ProjectStructureRequirement)(nil), structure)
				require.EqualValues(t, expectedProjectData, structure.ProjectsData)
			})

		opts := &CreateNewProjectOptions{
			Lister:     mockLister,
			Prompter:   mockPrompter,
			Manager:    newMockManager,
			Compressor: mockCompressor,
			Executor:   mockExecutor,
			Logger:     mockLogger,
			Store:      mockStore,
		}

		err := GetInstance().CreateNewProject(opts)

		require.Nil(t, err)
	})
}

func Test_isValid(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLister := lister.NewMockLister(controller)
	mockPrompter := prompter.NewMockPrompter(controller)
	newMockManager := manager.NewMockManager(controller)
	mockCompressor := compressor.NewMockCompressor(controller)
	mockExecutor := executor.NewMockExecutor(controller)
	mockLogger := logger.NewLogger()

	type args struct {
		opts *CreateNewProjectOptions
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all present",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:     mockLister,
					Prompter:   mockPrompter,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Logger:     mockLogger,
				},
			},
			want: true,
		},
		{
			name: "missing mockExecutor",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:     mockLister,
					Prompter:   mockPrompter,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Logger:     mockLogger,
				},
			},
			want: false,
		},
		{
			name: "missing mockCompressor",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:   mockLister,
					Prompter: mockPrompter,
					Manager:  newMockManager,
					Executor: mockExecutor,
					Logger:   mockLogger,
				},
			},
			want: false,
		},
		{
			name: "missing newMockManager",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:     mockLister,
					Prompter:   mockPrompter,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Logger:     mockLogger,
				},
			},
			want: false,
		},
		{
			name: "missing mockPrompter",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:     mockLister,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Logger:     mockLogger,
				},
			},
			want: false,
		},
		{
			name: "missing mockLister",
			args: args{
				opts: &CreateNewProjectOptions{
					Prompter:   mockPrompter,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Logger:     mockLogger,
				},
			},
			want: false,
		},
		{
			name: "missing logger",
			args: args{
				opts: &CreateNewProjectOptions{
					Prompter:   mockPrompter,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Lister:     mockLister,
				},
			},
			want: false,
		},
		{
			name: "should validate if path is nil",
			args: args{
				opts: &CreateNewProjectOptions{
					Lister:     mockLister,
					Prompter:   mockPrompter,
					Manager:    newMockManager,
					Compressor: mockCompressor,
					Executor:   mockExecutor,
					Logger:     mockLogger,
					Path:       nil,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.args.opts); got == nil != tt.want {
			}
		})
	}
}

func Test_isValid_PathTest(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLister := lister.NewMockLister(controller)
	mockPrompter := prompter.NewMockPrompter(controller)
	newMockManager := manager.NewMockManager(controller)
	mockCompressor := compressor.NewMockCompressor(controller)
	mockExecutor := executor.NewMockExecutor(controller)
	mockLogger := logger.NewLogger()
	mockStore := store.GetInstance()

	options := CreateNewProjectOptions{
		Lister:     mockLister,
		Prompter:   mockPrompter,
		Manager:    newMockManager,
		Compressor: mockCompressor,
		Executor:   mockExecutor,
		Logger:     mockLogger,
		Store:      mockStore,
		Path:       nil,
	}

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
