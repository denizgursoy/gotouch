package operator

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/go-playground/validator/v10"
)

var yamlValidSuccessMsg = "YAML file is valid"

type (
	ValidateYamlOptions struct {
		Lister lister.Lister `validate:"required"`
		Logger logger.Logger `validate:"required"`
		Path   *string       `validate:"required,endswith=.yaml,url|file"`
	}
)

func (o *operator) ValidateYaml(opts *ValidateYamlOptions) error {
	if validationError := isValidYaml(opts); validationError != nil {
		return validationError
	}

	_, err := opts.Lister.GetProjectList(opts.Path)
	if err != nil {
		return err
	}

	opts.Logger.LogInfo(yamlValidSuccessMsg)

	return nil
}

func isValidYaml(opts *ValidateYamlOptions) error {
	err := validator.New().Struct(opts)
	if err != nil {
		fieldErrors := err.(validator.ValidationErrors)
		fieldError := fieldErrors[0]
		if fieldError.Field() == "Path" {
			if fieldError.ActualTag() == "endswith" {
				return ErrNotYamlFile
			}
			return ErrNotValidUrlOrFilePath
		}

		return ErrAllFieldsAreRequired
	}
	return nil
}
