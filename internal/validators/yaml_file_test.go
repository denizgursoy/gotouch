package validators_test

import (
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/denizgursoy/gotouch/internal/validators"
)

func Test_yamlFileValidator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		args      any
		wantError bool
	}{
		{
			name: "Empty pointer",
			args: struct {
				Path *string `validate:"yaml_file"`
			}{},
			wantError: false,
		},
		{
			name: "valid YAML file",
			args: struct {
				Path string `validate:"yaml_file"`
			}{
				Path: filepath.Join("testdata", "valid.yml"),
			},
			wantError: false,
		},
		{
			name: "Valid pointer to YAML file",
			args: struct {
				Path *string `validate:"yaml_file"`
			}{
				Path: ptr(filepath.Join("testdata", "valid.yml")),
			},
			wantError: false,
		},
		{
			name: "Wrong path",
			args: struct {
				Path string `validate:"yaml_file"`
			}{
				Path: filepath.Join("testdata", "valid.yaml"),
			},
			wantError: true,
		},
		{
			name: "invalid YAML file",
			args: struct {
				Path string `validate:"yaml_file"`
			}{
				Path: filepath.Join("testdata", "invalid.yml"),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			validate := validator.New()
			err := validators.AddYamlFileValidator(validate)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			err = validate.Struct(tt.args)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
