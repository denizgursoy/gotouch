// +build unit

package executor

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	testRequirement struct {
		isAskCalled  bool
		Error        error
		Tasks        []model.Task
		Requirements []model.Requirement
	}

	testTask struct {
		isAskCalled bool
		returnValue interface{}
		err         error
	}
)

func (t *testTask) Complete() error {
	t.isAskCalled = true
	if t.err != nil {
		return t.err
	}
	return nil
}

func (t *testRequirement) AskForInput() ([]model.Task, []model.Requirement, error) {
	t.isAskCalled = true
	if t.Error != nil {
		return nil, nil, t.Error
	}
	return t.Tasks, t.Requirements, nil
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

		require.True(t, firstRequirement.Tasks[0].(*testTask).isAskCalled)
		require.True(t, secondRequirement.Tasks[0].(*testTask).isAskCalled)
	})

	t.Run("should return error if complete function of task returns error", func(t *testing.T) {
		executor := newExecutor()
		errorRequirement, taskError := getCompleteErrorRequirement()
		err := executor.Execute(Requirements{errorRequirement})

		require.True(t, errorRequirement.isAskCalled)
		require.True(t, errorRequirement.Tasks[0].(*testTask).isAskCalled)

		require.NotNil(t, err)
		require.ErrorIs(t, taskError, err)
	})

	t.Run("should return error if complete function of task returns error", func(t *testing.T) {
		executor := newExecutor()
		req := getRequirementReturningTwoRequirements()
		err := executor.Execute(Requirements{req})

		require.True(t, req.Requirements[0].(*testRequirement).isAskCalled)
		require.True(t, req.Requirements[1].(*testRequirement).isAskCalled)

		require.True(t, req.Requirements[0].(*testRequirement).Tasks[0].(*testTask).isAskCalled)
		require.True(t, req.Requirements[1].(*testRequirement).Tasks[0].(*testTask).isAskCalled)

		require.True(t, req.isAskCalled)
		require.True(t, req.Tasks[0].(*testTask).isAskCalled)

		require.NoError(t, err)
	})
}

func getRequirements() (Requirements, *testRequirement, *testRequirement) {
	requirements := make(Requirements, 0)

	firstRequirement := &testRequirement{
		isAskCalled: false,
		Error:       nil,
		Tasks: []model.Task{
			&testTask{
				isAskCalled: false,
				returnValue: "test return value",
			},
		},
	}

	secondRequirement := &testRequirement{
		isAskCalled: false,
		Error:       nil,
		Tasks: []model.Task{
			&testTask{
				isAskCalled: false,
				returnValue: nil,
			},
		},
	}

	requirements = append(requirements, firstRequirement, secondRequirement)
	return requirements, firstRequirement, secondRequirement
}

func getCompleteErrorRequirement() (*testRequirement, error) {
	completeError := errors.New("could not complete the test")
	errorRequirement := &testRequirement{
		isAskCalled: false,
		Error:       nil,
		Tasks: []model.Task{
			&testTask{
				isAskCalled: false,
				returnValue: nil,
				err:         completeError,
			},
		},
	}

	return errorRequirement, completeError
}

func getRequirementReturningTwoRequirements() *testRequirement {
	requirements, _, _ := getRequirements()
	return &testRequirement{
		isAskCalled:  false,
		Error:        nil,
		Requirements: requirements,
		Tasks: []model.Task{
			&testTask{
				isAskCalled: false,
				returnValue: nil,
				err:         nil,
			},
		},
	}

}
