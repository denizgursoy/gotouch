package requirements

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/cloner"
	"github.com/denizgursoy/gotouch/internal/commandrunner"
	"github.com/denizgursoy/gotouch/internal/langs"
	"testing"

	"github.com/denizgursoy/gotouch/internal/compressor"
	"github.com/denizgursoy/gotouch/internal/executor"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompter"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
		Name:      "Project -1",
		Reference: "go.dev",
		URL:       "https://project1.com",
		Questions: questions,
		Values: map[interface{}]interface{}{
			1: "23",
		},
	}
	projectStructure2 = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
	}

	projectStructureWithValues = model.ProjectStructureData{
		Name:      "Project -2",
		Reference: "go2.dev",
		URL:       "https://project2.com",
		Values: map[interface{}]interface{}{
			"x": "",
		},
	}

	testProjectData = []*model.ProjectStructureData{&projectStructure1, &projectStructure2}
)

func TestStructure_AskForInput(t *testing.T) {
	t.Run("should ask for selection for project", func(t *testing.T) {
		requirement, controller := getTestProjectRequirement(t, testProjectData)
		defer controller.Finish()

		options := make([]fmt.Stringer, 0)
		for _, datum := range testProjectData {
			options = append(options, datum)
		}

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForSelectionFromList(gomock.Eq(SelectProjectTypeDirection), gomock.Eq(options)).
			Return(testProjectData[0], nil).
			Times(1)

		requirement.Prompter.(*prompter.MockPrompter).
			EXPECT().
			AskForString(gomock.Eq(ModuleNameDirection), gomock.Any()).
			Return("", nil).
			Times(1)

		tasks, requirements, err := requirement.AskForInput()

		require.NoError(t, err)
		require.NotNil(t, tasks)

		require.Len(t, tasks, 2)
		require.IsType(t, (*projectNameTask)(nil), tasks[0])

		require.IsType(t, &projectStructureTask{}, tasks[1])
		require.IsType(t, testProjectData[1], tasks[1].(*projectStructureTask).ProjectStructure)

		require.IsType(t, (*cleanupTask)(nil), tasks[2])

		actualQuestions := make([]*model.Question, 0)

		require.Len(t, requirements, 5)
		for i := 0; i < 2; i++ {
			actualQuestions = append(actualQuestions, &requirements[i].(*QuestionRequirement).Question)
		}
		require.Equal(t, questions, actualQuestions)

		require.IsType(t, &templateRequirement{}, requirements[2])
		require.IsType(t, testProjectData[0].Values, requirements[2].(*templateRequirement).Values)
		require.IsType(t, (*cleanupRequirement)(nil), requirements[3])
		require.IsType(t, (*initRequirement)(nil), requirements[4])
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
		task, controller := getTestProjectTask(t)
		defer controller.Finish()

		task.Compressor.(*compressor.MockCompressor).
			EXPECT().
			UncompressFromUrl(gomock.Eq(projectStructure1.URL)).
			Return(nil)

		task.LanguageChecker.(*langs.MockChecker).EXPECT().Setup().Times(1)

		err := task.Complete()
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
	mockStore := store.GetInstance()
	mockCloner := cloner.NewMockCloner(controller)
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
		CommandRunner:   mockRunner,
	}, controller

}

func getTestProjectTask(t *testing.T) (projectStructureTask, *gomock.Controller) {
	controller := gomock.NewController(t)

	mockUncompressor := compressor.NewMockCompressor(controller)
	mockManager := manager.NewMockManager(controller)
	mockLogger := logger.NewLogger()
	mockExecutor := executor.NewMockExecutor(controller)
	mockStore := store.NewMockStore(controller)
	mockChecker := langs.NewMockChecker(controller)
	mockCloner := cloner.NewMockCloner(controller)

	return projectStructureTask{
		ProjectStructure: &projectStructure1,
		Compressor:       mockUncompressor,
		Manager:          mockManager,
		Logger:           mockLogger,
		Executor:         mockExecutor,
		Store:            mockStore,
		LanguageChecker:  mockChecker,
		Cloner:           mockCloner,
	}, controller
}
