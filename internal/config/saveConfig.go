package config

import (
	"encoding/json"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/logger"
	"os"
	"path/filepath"
	"strings"
)

type (
	configManager struct {
		logger logger.Logger
	}
	Config struct {
		DefaultPath *string `json:"default_path"`
	}
)

const (
	URL = "url"
)

var (
	ConfigurableSettings = []string{URL}
)

func (c *configManager) SetValueOf(key, value string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	if strings.TrimSpace(key) == URL {
		config.DefaultPath = &value
	}

	err = saveConfig(config)
	if err != nil {
		return err
	}
	c.logger.LogInfo(fmt.Sprintf("%s was set to %s succesfully", key, value))

	return nil
}

func (c *configManager) UnsetValuesOf(key string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	if strings.TrimSpace(key) == URL {
		config.DefaultPath = nil
	}

	err = saveConfig(config)
	if err != nil {
		return err
	}
	c.logger.LogInfo(fmt.Sprintf("%s was unset succesfully", key))

	return nil
}

func (c *configManager) GetDefaultPath() (string, error) {
	config, err := readConfig()
	if err != nil {
		return "", err
	}
	if config == nil || config.DefaultPath == nil || len(strings.TrimSpace(*config.DefaultPath)) == 0 {
		return PropertiesUrlAddress, nil
	} else {
		return *config.DefaultPath, nil
	}
}

func readConfig() (*Config, error) {
	name, err := GetFileName()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(name)

	if err != nil {
		return &Config{}, nil
	}

	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetFileName() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

func saveConfig(config *Config) error {
	name, err := GetFileName()
	if err != nil {
		return err
	}

	err = createFileIfNotExists(err, name)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(name, marshal, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createFileIfNotExists(err error, name string) error {
	if _, err = os.Stat(name); os.IsNotExist(err) {
		_, err = os.Create(name)
		if err != nil {
			return err
		}
	}
	return nil
}
