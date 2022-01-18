package common

type Config struct {
	Dependencies      []Dependency       `yaml:"dependencies"`
	ProjectStructures []ProjectStructure `yaml:"project_structures"`
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
	err := ReadYaml(c, "default.yaml")
	if err != nil {
		return nil
	}
	return c
}
