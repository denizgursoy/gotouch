//go:generate mockgen -source=./list.go -destination=mock-list.go -package=lister

package lister

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type (
	ProjectStructureData struct {
		Name      string `yaml:"name"`
		Reference string `yaml:"reference"`
		URL       string `yaml:"url"`
	}

	Lister interface {
		GetDefaultProjects() ([]*ProjectStructureData, error)
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

func ParseToProjectStructureData(reader io.ReadCloser) ([]*ProjectStructureData, error) {
	data := make([]*ProjectStructureData, 0)

	allBytes, err := ioutil.ReadAll(reader)
	err = yaml.Unmarshal(allBytes, &data)

	if err != nil {
		log.Println(err.Error())
		return nil, ProjectDataParseError
	}

	if data == nil {
		return nil, NoProjectError
	}

	return data, nil
}

func (p *ProjectStructureData) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Reference)
}
