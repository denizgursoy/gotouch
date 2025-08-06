//go:generate mockgen -source=./config.go -destination=mockConfig.go -package=config --typed
package config

import "github.com/denizgursoy/gotouch/internal/logger"

type (
	ConfigManager interface {
		SetValueOf(key, value string) error
		UnsetValuesOf(key string) error
		GetDefaultPath() (string, error)
	}
)

const (
	ConfigFileName       = ".gotouch"
	PropertiesUrlAddress = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml"
)

func NewConfigManager(logger logger.Logger) ConfigManager {
	return &configManager{
		logger: logger,
	}
}
