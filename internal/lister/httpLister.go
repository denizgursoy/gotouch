package lister

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/model"
	"net/http"
)

var (
	ConnectionError       = errors.New("could not fetch project from remote server")
	ProjectDataParseError = errors.New("data could not be parsed properly")
	NoProjectError        = errors.New("data cannot be empty")
	PropertiesUrl         = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/properties.yaml"
)

type httpLister struct {
	client *http.Client
	URL    *string
}

func newHttpLister(client *http.Client, url *string) Lister {
	return httpLister{
		client: client,
		URL:    url,
	}
}

func (h httpLister) GetDefaultProjects() ([]*model.ProjectStructureData, error) {
	response, err := h.client.Get(*h.URL)

	if err != nil {
		return nil, ConnectionError
	}

	return ParseToProjectStructureData(response.Body)
}
