package root

import (
	"testing"
)

func TestGetPackageCommandHandler(t *testing.T) {
	//t.Run("should package successfully", func(t *testing.T) {
	//
	//	type arg struct {
	//		flag    string
	//		pointer *string
	//	}
	//
	//	flag := "./test-input.yaml"
	//	arguments := []arg{
	//		{
	//			flag:    flag,
	//			pointer: &flag,
	//		},
	//		{
	//			flag:    "",
	//			pointer: nil,
	//		},
	//	}
	//
	//	for _, argument := range arguments {
	//		controller := gomock.NewController(t)
	//		mockCommander := operator.NewMockCommander(controller)
	//
	//		expectedCall := &operator.CreateCommandOptions{
	//			Lister:     lister.GetInstance(),
	//			Prompter:   prompter.GetInstance(),
	//			Manager:    manager.GetInstance(),
	//			Compressor: compressor.GetInstance(),
	//			Executor:   executor.GetInstance(),
	//			Logger:     logger.NewLogger(),
	//			Path:       argument.pointer,
	//			Store:      store.GetInstance(),
	//		}
	//
	//		mockCommander.EXPECT().CreateNewProject(gomock.Eq(expectedCall))
	//
	//		command := CreateRootCommand(mockCommander)
	//		command.SetArgs(getPackageTestArguments("".""))
	//
	//		err := command.Execute()
	//		require.Nil(t, err)
	//	}
	//
	//})
}

func getPackageTestArguments(source, target string) []string {
	return []string{"-s", source, "-t", target}
}
