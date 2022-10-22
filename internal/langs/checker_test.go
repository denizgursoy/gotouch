package langs

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLanguageChecker(t *testing.T) {
	t.Run("should be golang setup if language is go/golang or empty", func(t *testing.T) {
		strings := make([]string, 0)
		for _, s := range strings {
			checker := GetInstance()
			checker.Init(s, nil, nil)
			require.IsType(t, checker.GetLangChecker(), &golangSetupChecker{})
		}
	})

	t.Run("should be empty setup checker if language is not go/golang or empty ", func(t *testing.T) {
		checker := GetInstance()
		checker.Init("test", nil, nil)
		require.IsType(t, checker.GetLangChecker(), &emptySetupChecker{})
	})
}
