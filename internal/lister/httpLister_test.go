// +build unit

package lister

import (
	"bytes"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

var (
	testProjectStructures = make([]*model.ProjectStructureData, 0)

	project1 = model.ProjectStructureData{
		Name:      "sds",
		Reference: "sds",
		URL:       "sds",
	}

	project2 = model.ProjectStructureData{
		Name:      "Project-1",
		Reference: "sds",
		URL:       "sds",
	}
)

func init() {
	testProjectStructures = append(testProjectStructures, &project1, &project2)
}

func Test_gitHubLister_GetDefaultProjects(t *testing.T) {
	t.Run("should send request", func(t *testing.T) {

		client := NewTestClient(func(req *http.Request) *http.Response {
			marshal, _ := yaml.Marshal(testProjectStructures)
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(marshal)),
			}
		})

		githubLister := newHttpLister(client, &PropertiesUrl)

		projects, err := githubLister.GetDefaultProjects()

		require.NoError(t, err)
		require.NotNil(t, projects)

		require.EqualValues(t, projects, testProjectStructures)

	})

	t.Run("should return ConnectionError if cannot reach to remote server", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return nil
		})
		githubLister := newHttpLister(client, &PropertiesUrl)

		projects, err := githubLister.GetDefaultProjects()

		require.Nil(t, projects)
		require.NotNil(t, err)

		require.ErrorIs(t, err, ConnectionError)

	})

	t.Run("should return ProjectDataParseError if body cannot be parsed properly", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte("parse-error-only-string"))),
			}
		})
		githubLister := newHttpLister(client, &PropertiesUrl)

		projects, err := githubLister.GetDefaultProjects()

		require.Nil(t, projects)
		require.NotNil(t, err)

		require.ErrorIs(t, err, ProjectDataParseError)

	})

	t.Run("should return ProjectDataParseError if body cannot be parsed properly", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			marshal, _ := yaml.Marshal(nil)
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(marshal)),
			}
		})

		githubLister := newHttpLister(client, &PropertiesUrl)

		projects, err := githubLister.GetDefaultProjects()

		require.Nil(t, projects)
		require.NotNil(t, err)

		require.ErrorIs(t, err, NoProjectError)

	})

}
