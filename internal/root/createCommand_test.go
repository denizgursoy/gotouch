// +build unit

package root

import (
	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
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
		mockPrompter := prompter.NewMockPrompter(controller)
		newMockManager := manager.NewMockManager(controller)
		mockCompressor := compressor.NewMockCompressor(controller)
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

				require.NotNil(t, name.Prompter)
				require.NotNil(t, name.Manager)

				require.NotNil(t, structure.Compressor)
				require.NotNil(t, structure.Manager)
				require.NotNil(t, structure.Prompter)

				require.IsType(t, (*req.ProjectNameRequirement)(nil), name)
				require.IsType(t, (*req.ProjectStructureRequirement)(nil), structure)
				require.EqualValues(t, expectedProjectData, structure.ProjectsData)
			})

		opts := &CreateCommandOptions{
			lister:     mockLister,
			prompter:   mockPrompter,
			manager:    newMockManager,
			compressor: mockCompressor,
			executor:   mockExecutor,
		}
		err := CreateNewProject(opts)

		require.Nil(t, err)
	})
}

func Test_isValid(t *testing.T) {
	controller := gomock.NewController(t)
	mockLister := lister.NewMockLister(controller)
	mockPrompter := prompter.NewMockPrompter(controller)
	newMockManager := manager.NewMockManager(controller)
	mockCompressor := compressor.NewMockCompressor(controller)
	mockExecutor := executor.NewMockExecutor(controller)

	type args struct {
		opts *CreateCommandOptions
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all present",
			args: args{
				opts: &CreateCommandOptions{
					lister:     mockLister,
					prompter:   mockPrompter,
					manager:    newMockManager,
					compressor: mockCompressor,
					executor:   mockExecutor,
				},
			},
			want: true,
		},
		{
			name: "missing mockExecutor",
			args: args{
				opts: &CreateCommandOptions{
					lister:     mockLister,
					prompter:   mockPrompter,
					manager:    newMockManager,
					compressor: mockCompressor,
				},
			},
			want: false,
		},
		{
			name: "missing mockCompressor",
			args: args{
				opts: &CreateCommandOptions{
					lister:   mockLister,
					prompter: mockPrompter,
					manager:  newMockManager,
					executor: mockExecutor,
				},
			},
			want: false,
		},
		{
			name: "missing newMockManager",
			args: args{
				opts: &CreateCommandOptions{
					lister:     mockLister,
					prompter:   mockPrompter,
					compressor: mockCompressor,
					executor:   mockExecutor,
				},
			},
			want: false,
		},
		{
			name: "missing mockPrompter",
			args: args{
				opts: &CreateCommandOptions{
					lister:     mockLister,
					manager:    newMockManager,
					compressor: mockCompressor,
					executor:   mockExecutor,
				},
			},
			want: false,
		},
		{
			name: "missing mockLister",
			args: args{
				opts: &CreateCommandOptions{
					prompter:   mockPrompter,
					manager:    newMockManager,
					compressor: mockCompressor,
					executor:   mockExecutor,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.args.opts); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
