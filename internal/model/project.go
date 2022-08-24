package model

import "fmt"

type (
	ProjectStructureData struct {
		Name      string `yaml:"name"`
		Reference string `yaml:"reference"`
		URL       string `yaml:"url"`
	}
)

func (p *ProjectStructureData) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Reference)
}
