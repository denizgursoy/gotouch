//go:generate mockgen -source=./lister.go -destination=mockLister.go -package=lister

package lister

import (
	"errors"
	"github.com/denizgursoy/gotouch/internal/model"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type (
	Lister interface {
		GetProjectList(path *string) ([]*model.ProjectStructureData, error)
	}

	ReadStrategy interface {
		ReadProjectStructures() (io.ReadCloser, error)
	}

	mainLister struct {
		defaultStrategy ReadStrategy
	}
)

var (
	lister                Lister
	once                  sync.Once
	PropertiesUrl         = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/properties.yaml"
	ProjectDataParseError = errors.New("data could not be parsed properly")
	NoProjectError        = errors.New("data cannot be empty")
)

func GetInstance() Lister {
	once.Do(func() {
		uri, _ := url.ParseRequestURI(PropertiesUrl)
		lister = &mainLister{
			defaultStrategy: NewUrlReader(uri, &http.Client{}),
		}
	})
	return lister
}

func (m *mainLister) GetProjectList(path *string) ([]*model.ProjectStructureData, error) {
	strategy := m.defaultStrategy

	if path != nil && len(strings.TrimSpace(*path)) != 0 {
		strategy = determineReadStrategy(*path)
	}

	return m.getProjectsFromStrategy(strategy)
}

func (m *mainLister) getProjectsFromStrategy(strategy ReadStrategy) ([]*model.ProjectStructureData, error) {
	structures, readError := strategy.ReadProjectStructures()

	if readError != nil {
		return nil, readError
	}

	data, err := ParseToProjectStructureData(structures)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseToProjectStructureData(reader io.ReadCloser) ([]*model.ProjectStructureData, error) {
	data := make([]*model.ProjectStructureData, 0)

	allBytes, err := ioutil.ReadAll(reader)
	err = yaml.Unmarshal(allBytes, &data)

	if err != nil {
		return nil, ProjectDataParseError
	}

	if data == nil {
		return nil, NoProjectError
	}

	return data, nil
}

func determineReadStrategy(path string) ReadStrategy {
	uri, err := url.ParseRequestURI(path)
	if err != nil {
		return NewFileReader(path)
	} else {
		return NewUrlReader(uri, &http.Client{})
	}
}
