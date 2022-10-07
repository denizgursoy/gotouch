//go:build unit
// +build unit

package lister

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

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

type RoundTripFunction func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunction) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunction) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func init() {
	testProjectStructures = append(testProjectStructures, &project1, &project2)
}

func Test_urlReader_ReadProjectStructures(t *testing.T) {
	t.Run("should return successfully", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			marshal, _ := yaml.Marshal(testProjectStructures)
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(marshal)),
			}
		})
		parse, _ := url.Parse(PropertiesUrl)
		reader := NewUrlReader(parse, client)
		structures, err := reader.ReadProjectStructures()

		require.NotNil(t, structures)
		require.Nil(t, err)
	})

	t.Run("should return ConnectionError if cannot reach to remote server", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return nil
		})
		parse, _ := url.Parse(PropertiesUrl)
		reader := NewUrlReader(parse, client)
		structures, err := reader.ReadProjectStructures()

		require.Nil(t, structures)
		require.NotNil(t, err)

		require.ErrorIs(t, err, ConnectionError)
	})
}
