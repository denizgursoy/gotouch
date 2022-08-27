// +build unit

package req

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testProjectName = "test-project"
	testUrlName     = "github.com/user/test-project"
	extractLocation = "/tmp/var"
)

func Test_validateProjectName(t *testing.T) {
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success test 1",
			args:    args{projectName: "github.com/test/project"},
			wantErr: false,
		},
		{
			name:    "success test 2",
			args:    args{projectName: "github.com/test.com/project"},
			wantErr: false,
		},
		{
			name:    "success test 3",
			args:    args{projectName: "github123.com/test123.com/project123"},
			wantErr: false,
		},
		{
			name:    "success test 4",
			args:    args{projectName: "github123.com/test123.com/project123/project"},
			wantErr: false,
		},
		{
			name:    "success test 5",
			args:    args{projectName: "github123"},
			wantErr: false,
		},
		{
			name:    "success test 6",
			args:    args{projectName: "github"},
			wantErr: false,
		},
		{
			name:    "error test 1",
			args:    args{projectName: ""},
			wantErr: true,
		},
		{
			name:    "error test 2",
			args:    args{projectName: "."},
			wantErr: true,
		},
		{
			name:    "error test 3",
			args:    args{projectName: ".exe"},
			wantErr: true,
		},
		{
			name:    "error test 4",
			args:    args{projectName: "./test"},
			wantErr: true,
		},
		{
			name:    "error test 5",
			args:    args{projectName: "123test"},
			wantErr: true,
		},
		{
			name:    "error test 6",
			args:    args{projectName: "error.com/123"},
			wantErr: true,
		},
		{
			name:    "error test 7",
			args:    args{projectName: "error.com/test123/."},
			wantErr: true,
		},
		{
			name:    "error test 8",
			args:    args{projectName: "error.com/test123/blabla."},
			wantErr: true,
		},
		{
			name:    "error test 9",
			args:    args{projectName: "error.com/test123/blabla.exe"},
			wantErr: true,
		},
		{
			name:    "error test 10",
			args:    args{projectName: "error.com/test123.com"},
			wantErr: true,
		},
		{
			name:    "error test 11",
			args:    args{projectName: "error.111/test"},
			wantErr: true,
		},
		{
			name:    "error test 12",
			args:    args{projectName: "error.111"},
			wantErr: true,
		},
		{
			name:    "error test 13",
			args:    args{projectName: "error/errr./test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateProjectName(tt.args.projectName); (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectNameRequirement_AskForInput(t *testing.T) {
	t.Run("should operate successfully", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockLogger := logger.NewLogger()

		mockPrompter.
			EXPECT().
			AskForString(gomock.Any(), gomock.Any()).
			Return(testProjectName, nil).
			Times(1)

		requirement := ProjectNameRequirement{
			mockPrompter,
			mockManager,
			mockLogger,
		}

		input, err := requirement.AskForInput()
		if err != nil {
			return
		}

		require.NoError(t, err)
		require.NotNil(t, input)

		task := input.(*projectNameTask)
		require.NotNil(t, task.Manager)
		require.EqualValues(t, testProjectName, task.ProjectName)
		require.NotNil(t, task.Manager)
	})

	t.Run("should return error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockLogger := logger.NewLogger()

		inputErr := errors.New("input error")
		mockPrompter.
			EXPECT().
			AskForString(gomock.Any(), gomock.Any()).
			Return("", inputErr).
			Times(1)

		requirement := ProjectNameRequirement{
			mockPrompter,
			mockManager,
			mockLogger,
		}

		input, err := requirement.AskForInput()

		require.NotNil(t, err)
		require.Nil(t, input)
		require.ErrorIs(t, inputErr, err)
	})

}

func Test_projectNameTask_Complete(t *testing.T) {
	t.Run("should create directories successfully", func(t *testing.T) {
		type args struct {
			projectName      string
			projectDirectory string
		}
		testCases := []args{
			{projectName: testProjectName, projectDirectory: extractLocation + "/" + testProjectName},
			{projectName: testUrlName, projectDirectory: extractLocation + "/" + testProjectName},
		}

		controller := gomock.NewController(t)
		defer controller.Finish()

		for _, testCase := range testCases {

			mockManager := manager.NewMockManager(controller)
			mockLogger := logger.NewLogger()

			mockManager.
				EXPECT().
				GetExtractLocation().
				Return(extractLocation).
				Times(1)

			mockManager.
				EXPECT().
				CreateDirectoryIfNotExists(gomock.Eq(testCase.projectDirectory))

			task := projectNameTask{
				ProjectName: testCase.projectName,
				Manager:     mockManager,
				Logger:      mockLogger,
			}

			complete, err := task.Complete(nil)
			require.NoError(t, err)
			require.NotNil(t, complete)

			actualOutput := complete.(string)

			require.EqualValues(t, testCase.projectName, actualOutput)
		}
	})
	t.Run("should return error if directory exists", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockManager := manager.NewMockManager(controller)

		mockManager.
			EXPECT().
			GetExtractLocation().
			Return(extractLocation).
			Times(1)

		mockManager.
			EXPECT().
			CreateDirectoryIfNotExists(gomock.Any()).
			Return(errors.New("could not create folder")).
			Times(1)

		task := projectNameTask{
			ProjectName: testProjectName,
			Manager:     mockManager,
		}

		complete, err := task.Complete(nil)

		require.NotNil(t, err)
		require.Nil(t, complete)
	})
}
