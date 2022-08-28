// +build unit

package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectStructureData_IsValid(t *testing.T) {
	projectName := "test-Project"
	validProjectURL := "https://github.com/denizgursoy/gotouch/graphs/traffic"

	validProjectWithNoCustom := ProjectStructureData{
		Name:      projectName,
		Reference: "",
		URL:       validProjectURL,
	}

	t.Run("should return ErrProjectNameIsEmpty if url is empty", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      "",
			Reference: "",
			URL:       "",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, ErrProjectNameIsEmpty, err)
	})

	t.Run("should return ErrProjectURLIsEmpty if url is empty", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      projectName,
			Reference: "",
			URL:       "",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, ErrProjectURLIsEmpty, err)
	})

	t.Run("should return ErrProjectURLIsNotValid if url is not valid", func(t *testing.T) {
		project := &ProjectStructureData{
			Name:      projectName,
			Reference: "",
			URL:       "test-url",
		}
		err := project.IsValid()

		require.NotNil(t, err)
		require.ErrorIs(t, ErrProjectURLIsNotValid, err)
	})

	t.Run("should return no error", func(t *testing.T) {
		err := validProjectWithNoCustom.IsValid()
		require.Nil(t, err)
	})

}
