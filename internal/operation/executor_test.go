// +build unit

package operation

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	testRequirement struct {
		isAskCalled bool
		Error       error
		Task        model.Task
	}

	testTask struct {
		isAskCalled bool
		arg         interface{}
		returnValue interface{}
	}
)

func (t *testTask) Complete(i interface{}) (interface{}, error) {
	t.isAskCalled = true
	t.arg = i
	return t.returnValue, nil
}

func (t *testRequirement) AskForInput() (model.Task, error) {
	t.isAskCalled = true
	if t.Error != nil {
		return nil, t.Error
	}
	return t.Task, nil
}

func Test_executor_Execute(t *testing.T) {
	t.Run("should return error if the requirement is nil", func(t *testing.T) {
		executor := newExecutor()
		err := executor.Execute(nil)
		require.ErrorIs(t, err, EmptyRequirementError)
	})

	t.Run("should call Ask For Input for all Requirements and complete the tasks", func(t *testing.T) {
		executor := newExecutor()
		requirements, firstRequirement, secondRequirement := getRequirements()

		err := executor.Execute(requirements)

		require.NoError(t, err)

		require.True(t, firstRequirement.isAskCalled)
		require.True(t, secondRequirement.isAskCalled)

		require.True(t, firstRequirement.Task.(*testTask).isAskCalled)
		require.True(t, secondRequirement.Task.(*testTask).isAskCalled)

		require.EqualValues(t, firstRequirement.Task.(*testTask).returnValue, secondRequirement.Task.(*testTask).arg)
	})
}

func getRequirements() (Requirements, *testRequirement, *testRequirement) {
	requirements := make(Requirements, 0)

	firstRequirement := &testRequirement{
		isAskCalled: false,
		Error:       nil,
		Task: &testTask{
			isAskCalled: false,
			arg:         nil,
			returnValue: "test return value",
		},
	}

	secondRequirement := &testRequirement{
		isAskCalled: false,
		Error:       nil,
		Task: &testTask{
			isAskCalled: false,
			arg:         nil,
			returnValue: nil,
		},
	}

	requirements = append(requirements, firstRequirement, secondRequirement)
	return requirements, firstRequirement, secondRequirement
}
