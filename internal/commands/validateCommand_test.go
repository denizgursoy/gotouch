package commands

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/operator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetValidateCommandHandler(t *testing.T) {
	t.Run("should validate successfully", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockCommander := operator.NewMockOperator(controller)

		path := "sdsd"
		expectedCall := &operator.ValidateYamlOptions{
			Lister: lister.GetInstance(),
			Logger: logger.NewLogger(),
			Path:   &path,
		}

		mockCommander.EXPECT().ValidateYaml(gomock.Eq(expectedCall))

		command := CreateValidateCommand(mockCommander)

		command.SetArgs(getValidateTestArguments(path))
		err := command.Execute()
		require.Nil(t, err)
	})
}

func getValidateTestArguments(source string) []string {
	return []string{"-f", source}
}
