package requirements

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
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
			wantErr: false,
		},
		{
			name:    "error test 4",
			args:    args{projectName: "./test"},
			wantErr: true,
		},
		{
			name:    "error test 5",
			args:    args{projectName: "123test"},
			wantErr: false,
		},
		{
			name:    "error test 6",
			args:    args{projectName: "error.com/123"},
			wantErr: false,
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
			wantErr: false,
		},
		{
			name:    "error test 10",
			args:    args{projectName: "error.com/test123.com"},
			wantErr: false,
		},
		{
			name:    "error test 11",
			args:    args{projectName: "error.111/test"},
			wantErr: false,
		},
		{
			name:    "error test 12",
			args:    args{projectName: "error.111"},
			wantErr: false,
		},
		{
			name:    "error test 13",
			args:    args{projectName: "error/errr./test"},
			wantErr: true,
		},
	}

	mockManager := manager.NewMockManager(gomock.NewController(t))
	mockManager.EXPECT().GetExtractLocation().AnyTimes()

	mockStore := store.NewMockStore(gomock.NewController(t))
	mockStore.EXPECT().GetValue(store.Inline).Return("false").AnyTimes()

	req := &ProjectNameRequirement{
		Manager: mockManager,
		Store:   mockStore,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := req.validateModuleName(tt.args.projectName); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleName() error = %v, wantErr %v, values %s", err, tt.wantErr, tt.args.projectName)
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
		mockStore := store.GetInstance()

		requirement := ProjectNameRequirement{
			mockPrompter,
			mockManager,
			mockLogger,
			mockStore,
			"test-initial-value",
		}

		mockPrompter.
			EXPECT().
			AskForString(gomock.Any(), requirement.InitialValue, gomock.Any()).
			Return(testProjectName, nil).
			Times(1)

		tasks, requirements, err := requirement.AskForInput()
		if err != nil {
			return
		}

		require.NoError(t, err)
		require.NotNil(t, tasks)
		require.Empty(t, requirements)

		task := tasks[0].(*projectNameTask)
		require.NotNil(t, task.Manager)
		require.EqualValues(t, testProjectName, task.ModuleName)
		require.NotNil(t, task.Manager)
	})

	t.Run("should return error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockPrompter := prompter.NewMockPrompter(controller)
		mockManager := manager.NewMockManager(controller)
		mockLogger := logger.NewLogger()
		mockStore := store.GetInstance()

		inputErr := errors.New("input error")
		requirement := ProjectNameRequirement{
			mockPrompter,
			mockManager,
			mockLogger,
			mockStore,
			"test-initial-value",
		}
		mockPrompter.
			EXPECT().
			AskForString(gomock.Any(), requirement.InitialValue, gomock.Any()).
			Return("", inputErr).
			Times(1)

		tasks, requirements, err := requirement.AskForInput()

		require.NotNil(t, err)
		require.Nil(t, tasks)
		require.Nil(t, requirements)
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
			mockStore := store.NewMockStore(controller)

			mockStore.EXPECT().SetValue(gomock.Any(), gomock.Any()).AnyTimes()
			mockStore.EXPECT().GetValue(store.Inline).Times(1).Return("false")

			mockManager.
				EXPECT().
				GetExtractLocation().
				Return(extractLocation).
				Times(1)

			mockManager.
				EXPECT().
				CreateDirectoryIfNotExist(gomock.Eq(testCase.projectDirectory))

			task := projectNameTask{
				ModuleName: testCase.projectName,
				Manager:    mockManager,
				Logger:     mockLogger,
				Store:      mockStore,
			}

			err := task.Complete(context.Background())
			require.NoError(t, err)
		}
	})
	t.Run("should return error if directory exists", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockManager := manager.NewMockManager(controller)
		mockLogger := logger.NewLogger()
		mockStore := store.NewMockStore(controller)

		mockStore.EXPECT().SetValue(gomock.Any(), gomock.Any()).AnyTimes()
		mockStore.EXPECT().GetValue(store.Inline).Times(1).Return("false")
		mockManager.
			EXPECT().
			GetExtractLocation().
			Return(extractLocation).
			Times(1)

		mockManager.
			EXPECT().
			CreateDirectoryIfNotExist(gomock.Any()).
			Return(errors.New("could not create folder")).
			Times(1)

		task := projectNameTask{
			ModuleName: testProjectName,
			Manager:    mockManager,
			Logger:     mockLogger,
			Store:      mockStore,
		}

		err := task.Complete(context.Background())

		require.NotNil(t, err)
	})
	t.Run("should not create directory if inline flag is present", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockManager := manager.NewMockManager(controller)
		mockLogger := logger.NewLogger()
		mockStore := store.NewMockStore(controller)

		mockStore.EXPECT().SetValue(store.ModuleName, testProjectName).Times(1)
		mockStore.EXPECT().SetValue(store.ProjectName, testProjectName).Times(1)
		mockStore.EXPECT().SetValue(store.WorkingDirectory, extractLocation).Times(1)
		mockStore.EXPECT().SetValue(store.ProjectFullPath, extractLocation).Times(1)
		mockStore.EXPECT().GetValue(store.Inline).Times(1).Return("true")

		mockManager.
			EXPECT().
			GetExtractLocation().
			Return(extractLocation).
			Times(1)

		task := projectNameTask{
			ModuleName: testProjectName,
			Manager:    mockManager,
			Logger:     mockLogger,
			Store:      mockStore,
		}

		err := task.Complete(context.Background())

		require.NoError(t, err)
	})
}
