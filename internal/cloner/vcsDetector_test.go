package cloner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_defaultVCSDetector_DetectVCS(t *testing.T) {
	t.Run("should return git it it is normal repository", func(t *testing.T) {
		detector := NewDefaultVCSDetector()
		vcs, err := detector.DetectVCS(t.Context(), "https://github.com/denizgursoy/gotouch")
		require.NoError(t, err)
		require.Equal(t, VCSGit, vcs)
	})
	t.Run("should return git it it is with ssh url", func(t *testing.T) {
		detector := NewDefaultVCSDetector()
		vcs, err := detector.DetectVCS(t.Context(), "git@github.com:denizgursoy/gotouch.git")
		require.NoError(t, err)
		require.Equal(t, VCSGit, vcs)
	})
	t.Run("should return none it it is git repository", func(t *testing.T) {
		detector := NewDefaultVCSDetector()
		vcs, err := detector.DetectVCS(t.Context(), "https://github.com/denizgursoy")
		require.NoError(t, err)
		require.Equal(t, VCSNone, vcs)
	})
}
