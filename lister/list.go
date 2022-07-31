package lister

import "sync"

type ProjectData struct {
	Name      string `yaml:"name"`
	Reference string `yaml:"reference"`
	URL       string `yaml:"url"`
}

type Lister interface {
	GetDefaultProjects() []*ProjectData
}

var (
	lister Lister
	once   sync.Once
)

func GetInstance() Lister {
	once.Do(func() {
		lister = newGithubLister()
	})
	return lister
}
