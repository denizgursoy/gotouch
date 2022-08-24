package model

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type (
	ProjectStructureData struct {
		Name      string      `yaml:"name"`
		Reference string      `yaml:"reference"`
		URL       string      `yaml:"url"`
		Questions []*Question `yaml:"Questions"`
	}

	Question struct {
		Direction         string    `yaml:"Direction"`
		CanSkip           bool      `yaml:"CanSkip"`
		CanSelectMultiple bool      `yaml:"CanSelectMultiple"`
		Options           []*Option `yaml:"options"`
	}

	Option struct {
		Answer       string    `yaml:"Answer"`
		Dependencies []*string `yaml:"Dependencies"`
		Files        []*File   `yaml:"Files"`
	}

	File struct {
		Url     string `yaml:"Url"`
		Content string `yaml:"Content"`
		Path    string `yaml:"Path"`
	}
)

func (p *ProjectStructureData) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Reference)
}

var (
	ErrProjectURLIsEmpty    = errors.New("project url can not be empty")
	ErrProjectNameIsEmpty   = errors.New("project name can not be empty")
	ErrProjectURLIsNotValid = errors.New("project name can not be empty")
)

func (p *ProjectStructureData) IsValid() error {
	if len(strings.TrimSpace(p.Name)) == 0 {
		return ErrProjectNameIsEmpty
	}

	projectUrl := strings.TrimSpace(p.URL)
	if len(projectUrl) == 0 {
		return ErrProjectURLIsEmpty
	}

	if _, err := url.ParseRequestURI(projectUrl); err != nil {
		return ErrProjectURLIsNotValid
	}

	return nil

}
