package common

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed default.yaml
var defaultFile string

type Config struct {
	Dependencies      []Dependency       `yaml:"dependencies"`
	ProjectStructures []ProjectStructure `yaml:"projectStructures"`
}

type Dependency struct {
	Prompt  string    `yaml:"prompt"`
	Options []Library `yaml:"options"`
}

type Library struct {
	Name    string
	Address string
	Version string
}

type ProjectStructure struct {
	Name        string   `yaml:"name"`
	Files       []File   `yaml:"files"`
	Directories []string `yaml:"directories"`
}

type File struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Content  string `yaml:"content"`
}

func GetDefaultConfig() *Config {
	c := &Config{}
	err := yaml.Unmarshal([]byte(defaultFile), c)
	if err != nil {
		return nil
	}
	return c
}
