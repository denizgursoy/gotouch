package operation

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_executor_Execute(t *testing.T) {
	t.Run("should return error if the requirement is nil", func(t *testing.T) {
		executor := newExecutor()
		err := executor.Execute(nil)
		require.ErrorIs(t, err, EmptyRequirementError)
	})
}
