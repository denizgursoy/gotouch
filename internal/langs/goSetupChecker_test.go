package langs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDisplayName(t *testing.T) {
	t.Run("should include reference if exists", func(t *testing.T) {
		version := "1234"
		url := "*"
		dependency := Dependency{Url: &url, Version: &version}
		require.EqualValues(t, url+"@"+version, dependency.String())
	})
}
