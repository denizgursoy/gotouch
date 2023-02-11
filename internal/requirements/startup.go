package requirements

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
)

type (
	startupTask struct {
		Store  store.Store   `validate:"required"`
		Logger logger.Logger `validate:"required"`
	}
)

const PropertiesYamlName = "properties.yaml"

func (s *startupTask) Complete() error {
	if err := validator.New().Struct(s); err != nil {
		return err
	}
	workingDirectory := s.Store.GetValue(store.ProjectFullPath)
	propertiesYamlFullAddress := filepath.Join(workingDirectory, PropertiesYamlName)
	if _, err := os.Stat(propertiesYamlFullAddress); !os.IsNotExist(err) {
		err := os.Remove(propertiesYamlFullAddress)
		if err != nil {
			s.Logger.LogErrorIfExists(fmt.Errorf("couuld not delete %s, error=%w", propertiesYamlFullAddress, err))
			return err
		}
		s.Logger.LogInfo(fmt.Sprintf("Deleted %s succesfully", PropertiesYamlName))
	}

	return nil
}
