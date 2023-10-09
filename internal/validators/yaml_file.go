package validators

import (
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func AddYamlFileValidator(validate *validator.Validate) error {
	return validate.RegisterValidation("yaml_file", yamlFileValidator, true)
}

func yamlFileValidator(fl validator.FieldLevel) bool {
	fieldVal := fl.Field()
	if fieldVal.Kind() == reflect.Pointer {
		if fieldVal.IsNil() {
			return true
		} else {
			fieldVal = fieldVal.Elem()
		}
	}

	if fieldVal.Kind() != reflect.String {
		return false
	}

	val := fieldVal.String()

	f, err := os.Open(val)
	if err != nil {
		return false
	}

	defer func() {
		_ = f.Close()
	}()

	return isYaml(f)
}
