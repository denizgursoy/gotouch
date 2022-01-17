package common

type Config struct {
	Hits int64 `yaml:"hits"`
	Time int64 `yaml:"time"`
}

func GetDefaultConfig() *Config {
	c := &Config{}
	err := ReadYaml(c, "default.yaml")
	if err != nil {
		return nil
	}
	return c
}
