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
		Questions []*Question `yaml:"questions"`
		Values    interface{} `yaml:"values"`
	}

	Question struct {
		Direction         string    `yaml:"direction"`
		CanSkip           bool      `yaml:"canSkip"`
		CanSelectMultiple bool      `yaml:"canSelectMultiple"`
		Options           []*Option `yaml:"options"`
	}

	Option struct {
		Answer       string    `yaml:"answer"`
		Dependencies []*string `yaml:"dependencies"`
		Files        []*File   `yaml:"files"`
	}

	File struct {
		Url     string `yaml:"url"`
		Content string `yaml:"content"`
		Path    string `yaml:"path"`
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

func (o *Option) String() string {
	return o.Answer
}

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
