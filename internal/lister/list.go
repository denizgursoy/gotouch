//go:generate mockgen -source=./list.go -destination=mock-list.go -package=lister

package lister

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type (
	Lister interface {
		GetDefaultProjects() ([]*model.ProjectStructureData, error)
	}
)

var (
	lister Lister
	once   sync.Once
)

func GetInstance() Lister {
	once.Do(func() {
		lister = newHttpLister(&http.Client{}, &PropertiesUrl)
	})
	return lister
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
