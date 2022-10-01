//go:generate mockgen -source=./lister.go -destination=mockLister.go -package=lister

package lister

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type (
	Lister interface {
		GetProjectList(path *string) ([]*model.ProjectStructureData, error)
	}

	ReadStrategy interface {
		ReadProjectStructures() (io.ReadCloser, error)
	}

	mainLister struct {
		DefaultStrategy ReadStrategy `validate:"required"`
	}
)

var (
	lister                Lister
	once                  sync.Once
	ProjectDataParseError = errors.New("data could not be parsed properly")
	NoProjectError        = errors.New("data cannot be empty")
)

func GetInstance() Lister {
	once.Do(func() {
		uri, _ := url.ParseRequestURI(PropertiesUrl)
		lister = &mainLister{
			DefaultStrategy: NewUrlReader(uri, &http.Client{}),
		}
	})
	return lister
}

func (m *mainLister) GetProjectList(path *string) ([]*model.ProjectStructureData, error) {
	if isValid(m) != nil {
		return nil, model.ErrMissingField
	}

	strategy := m.DefaultStrategy

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
	structures := make([]*model.ProjectStructureData, 0)

	allBytes, err := ioutil.ReadAll(reader)
	err = yaml.Unmarshal(allBytes, &structures)

	if err != nil {
		return nil, ProjectDataParseError
	}

	for _, structure := range structures {
		err = structure.IsValid()
		if err != nil {
			return nil, err
		}
	}

	if structures == nil {
		return nil, NoProjectError
	}

	return structures, nil
}

func determineReadStrategy(path string) ReadStrategy {
	_, err2 := os.Stat(path)
	if err2 != nil {
		uri, err := url.ParseRequestURI(path)
		if err != nil {
			fmt.Println(err)
		}
		return NewUrlReader(uri, &http.Client{})
	} else {
		return NewFileReader(path)
	}
}

func isValid(m *mainLister) error {
	return validator.New().Struct(m)
}
