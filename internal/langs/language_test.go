package langs

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLanguageChecker(t *testing.T) {
	t.Run("should be golang setup if language is go/golang or empty", func(t *testing.T) {
		strings := make([]string, 0)
		for _, s := range strings {
			data := model.ProjectStructureData{
				Language: s,
			}
			checker := NewLanguageChecker(&data)
			require.IsType(t, checker, &golangSetupChecker{})
		}
	})

	t.Run("should be empty setup checker if language is not go/golang or empty ", func(t *testing.T) {
		data := model.ProjectStructureData{
			Language: "test",
		}
		checker := NewLanguageChecker(&data)
		require.IsType(t, checker, &emptySetupChecker{})
	})
}
