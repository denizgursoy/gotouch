package requirements

import (
	"context"
	"fmt"
	"testing"

	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/langs"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
)

var (
	questions = []*model.Question{
		{
			Direction:         "question 1",
			CanSkip:           false,
			CanSelectMultiple: false,
			Choices:           nil,
		},
		{
			Direction:         "question 2",
			CanSkip:           false,
			CanSelectMultiple: false,
			Choices:           nil,
		},
	}
	projectStructure1 = model.ProjectStructureData{
		Name:              "Project -1",
		Reference:         "go.dev",
		URL:               "https://project1.com",
		InitialModuleName: "test-initial-module-name-1",
		Questions:         questions,
		Resources: model.Resources{
			Values: map[string]any{
				"1": "23",
			},
			CustomValues: map[string]any{
				"2": "23",
			},
			Dependencies: []any{"dep1", "dep2"},
			Files:        []*model.File{&file1, &file2},
		},
	}
	projectStructureWithGitRepository = model.ProjectStructureData{
		Name:              "Project -1",
		Reference:         "go.dev",
		URL:               "a.git",
		Branch:            "test",
		InitialModuleName: "test-initial-module-name-2",
		Questions:         questions,
		Resources: model.Resources{
			Values: map[string]any{
				"1": "23",
			},
		},
	}
	projectStructure2 = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
		Resources: model.Resources{
			Values: map[string]any{
				"x": "",
			},
		},
	}

	testProjectData = []*model.ProjectStructureData{
		&projectStructure1,
		&projectStructure2,
	}
)

func TestStructure_AskForInput(t *testing.T) {
	t.Run("should ask for selection for project", func(t *testing.T) {
		requirement, controller := getTestProjectRequirement(t, testProjectData)
		defer controller.Finish()

		options := make([]fmt.Stringer, 0)
		for _, datum := range testProjectData {
			options = append(options, datum)
		}

		selectedProjectStructure := testProjectData[0]
		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Eq(SelectProjectTypeDirection), gomock.Eq(options)).
			Return(selectedProjectStructure, nil).
			Times(1)

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForString(gomock.Eq(ModuleNameDirection), selectedProjectStructure.InitialModuleName, gomock.Any()).
			Return("", nil).
			Times(1)

		for _, dependency := range selectedProjectStructure.Dependencies {
			requirement.Store.(*store.MockStore).
				EXPECT().
				AddDependency(gomock.Eq(dependency))
		}

		requirement.Store.(*store.MockStore).
			EXPECT().
			AddCustomValues(gomock.Eq(selectedProjectStructure.CustomValues)).Times(1)

		requirement.Store.(*store.MockStore).
			EXPECT().
			AddValues(gomock.Eq(selectedProjectStructure.Values)).Times(1)

		tasks, requirements, err := requirement.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, tasks)

		require.Len(t, tasks, 7)
		require.IsType(t, (*projectNameTask)(nil), tasks[0])

		require.IsType(t, &projectStructureTask{}, tasks[1])
		require.Equal(t, selectedProjectStructure, tasks[1].(*projectStructureTask).ProjectStructure)

		require.IsType(t, (*dependencyTask)(nil), tasks[2])
		require.Equal(t, selectedProjectStructure.Resources.Dependencies[0], tasks[2].(*dependencyTask).Dependency)

		require.IsType(t, (*dependencyTask)(nil), tasks[3])
		require.Equal(t, selectedProjectStructure.Resources.Dependencies[1], tasks[3].(*dependencyTask).Dependency)

		require.IsType(t, (*fileTask)(nil), tasks[4])
		require.Equal(t, *selectedProjectStructure.Resources.Files[0], tasks[4].(*fileTask).File)

		require.IsType(t, (*fileTask)(nil), tasks[5])
		require.Equal(t, *selectedProjectStructure.Resources.Files[1], tasks[5].(*fileTask).File)

		require.IsType(t, (*startupTask)(nil), tasks[6])

		actualQuestions := make([]*model.Question, 0)

		require.Len(t, requirements, 5)
		for i := 0; i < 2; i++ {
			actualQuestions = append(actualQuestions, &requirements[i].(*QuestionRequirement).Question)
		}
		require.Equal(t, questions, actualQuestions)

		require.IsType(t, &templateRequirement{}, requirements[2])
		require.IsType(t, (*initRequirement)(nil), requirements[3])
		require.IsType(t, (*cleanupRequirement)(nil), requirements[4])
	})

	t.Run("should return error from the prompt", func(t *testing.T) {
		requirement, controller := getTestProjectRequirement(t, nil)
		defer controller.Finish()

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Any(), gomock.Any()).
			Return(nil, prompter.EmptyList).
			Times(1)

		tasks, _, err := requirement.AskForInput()

		require.NotNil(t, err)
		require.ErrorIs(t, err, prompter.EmptyList)

		require.Empty(t, tasks)
	})
}

func TestStructure_Complete(t *testing.T) {
	t.Run("should call uncompress with the URL", func(t *testing.T) {
		task := getTestProjectTask(t)

		task.VCSDetector.(*cloner.MockVCSDetector).
			EXPECT().
			DetectVCS(gomock.Any(), gomock.Eq(projectStructure1.URL)).
			Return(cloner.VCSNone, nil)

		task.Compressor.(*compressor.MockCompressor).
			EXPECT().
			UncompressFromUrl(gomock.Any(), gomock.Eq(projectStructure1.URL)).
			Return(nil)

		task.LanguageChecker.(*langs.MockChecker).EXPECT().Setup().Times(1)

		err := task.Complete(context.Background())
		require.Nil(t, err)
	})

	t.Run("should call cloner with branch name", func(t *testing.T) {
		task := getTestProjectTask(t)

		task.ProjectStructure = &projectStructureWithGitRepository
		task.VCSDetector.(*cloner.MockVCSDetector).
			EXPECT().
			DetectVCS(gomock.Any(), gomock.Eq(projectStructureWithGitRepository.URL)).
			Return(cloner.VCSGit, nil)
		task.Cloner.(*cloner.MockCloner).
			EXPECT().
			CloneFromUrl(gomock.Any(), gomock.Eq(projectStructureWithGitRepository.URL), gomock.Eq(projectStructureWithGitRepository.Branch)).
			Return(nil)

		task.LanguageChecker.(*langs.MockChecker).EXPECT().Setup().Times(1)

		err := task.Complete(context.Background())
		require.Nil(t, err)
	})
}

func getTestProjectRequirement(t *testing.T, projectData []*model.ProjectStructureData) (ProjectStructureRequirement, *gomock.Controller) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockPrompter := prompter.NewMockPrompter(controller)
	mockCompressor := compressor.NewMockCompressor(controller)
	mockManager := manager.NewMockManager(controller)
	mockExecutor := executor.NewMockExecutor(controller)
	mockLogger := logger.NewLogger()
	mockStore := store.NewMockStore(controller)
	mockCloner := cloner.NewMockCloner(controller)
	mockVCSDetector := cloner.NewMockVCSDetector(controller)
	mockRunner := commandrunner.NewMockRunner(controller)

	return ProjectStructureRequirement{
		ProjectsData:    projectData,
		Prompter:        mockPrompter,
		Compressor:      mockCompressor,
		Manager:         mockManager,
		Logger:          mockLogger,
		Executor:        mockExecutor,
		Store:           mockStore,
		LanguageChecker: langs.NewMockChecker(controller),
		Cloner:          mockCloner,
		VCSDetector:     mockVCSDetector,
		CommandRunner:   mockRunner,
	}, controller
}

func getTestProjectTask(t *testing.T) projectStructureTask {
	controller := gomock.NewController(t)

	mockUncompressor := compressor.NewMockCompressor(controller)
	mockManager := manager.NewMockManager(controller)
	mockLogger := logger.NewLogger()
	mockExecutor := executor.NewMockExecutor(controller)
	mockStore := store.NewMockStore(controller)
	mockChecker := langs.NewMockChecker(controller)
	mockCloner := cloner.NewMockCloner(controller)
	mockVCSDetector := cloner.NewMockVCSDetector(controller)

	return projectStructureTask{
		ProjectStructure: &projectStructure1,
		Compressor:       mockUncompressor,
		Manager:          mockManager,
		Logger:           mockLogger,
		Executor:         mockExecutor,
		Store:            mockStore,
		LanguageChecker:  mockChecker,
		Cloner:           mockCloner,
		VCSDetector:      mockVCSDetector,
	}
}
