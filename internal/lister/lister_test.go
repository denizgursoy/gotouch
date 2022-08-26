package lister

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestParseToProjectStructureData(t *testing.T) {
	t.Run("should return ProjectDataParseError if body cannot be parsed properly", func(t *testing.T) {
		closer := ioutil.NopCloser(bytes.NewReader([]byte("parse-error-only-string")))
		projects, err := ParseToProjectStructureData(closer)

		require.Nil(t, projects)
		require.NotNil(t, err)

		require.ErrorIs(t, err, ProjectDataParseError)
	})

	t.Run("should return NoProjectError if body cannot be parsed properly", func(t *testing.T) {
		marshal, _ := yaml.Marshal(nil)
		closer := ioutil.NopCloser(bytes.NewReader(marshal))

		projects, err := ParseToProjectStructureData(closer)

		require.Nil(t, projects)
		require.NotNil(t, err)

		require.ErrorIs(t, err, NoProjectError)
	})
}

func Test_determineReadStrategy(t *testing.T) {
	strategy := determineReadStrategy("")
	require.IsType(t, &urlReader{}, strategy)

	//TODO add file stragety test

	strategy = determineReadStrategy(PropertiesUrl)
	require.IsType(t, &urlReader{}, strategy)
}
