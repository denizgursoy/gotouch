package common

type Config struct {
	Dependencies []Dependency `yaml:"dependencies"`
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

func GetDefaultConfig() *Config {
	c := &Config{}
	err := ReadYaml(c, "default.yaml")
	if err != nil {
		return nil
	}
	return c
}
