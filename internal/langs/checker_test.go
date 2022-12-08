package langs

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLanguageChecker(t *testing.T) {
	t.Run("should be golang setup if language is go/golang or empty", func(t *testing.T) {
		strings := make([]string, 0)
		strings = append(strings, "go", "Go", "GO", "golang", "Golang", "GOLANG")
		for _, s := range strings {
			checker := GetChecker(s, nil, nil, nil)
			require.IsType(t, checker, &golangSetupChecker{})
		}
	})

	t.Run("should be empty setup checker if language is empty ", func(t *testing.T) {
		strings := make([]string, 0)
		strings = append(strings, " ", "    ", "", "test", "java", "go-test")
		for _, s := range strings {
			checker := GetChecker(s, nil, nil, nil)
			require.IsType(t, checker, &emptySetupChecker{})
		}
	})
}
