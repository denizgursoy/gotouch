package validators_test

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/denizgursoy/gotouch/internal/validators"
)

func Test_yamlUrlValidator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		args      any
		wantError bool
	}{
		{
			name: "Empty pointer",
			args: struct {
				Path *string `validate:"yaml_url"`
			}{},
			wantError: false,
		},
		{
			name: "valid YAML HTTP URL",
			args: struct {
				Path string `validate:"yaml_url"`
			}{
				Path: "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package.yaml",
			},
			wantError: false,
		},
		{
			name: "valid pointer to YAML HTTP URL",
			args: struct {
				Path *string `validate:"yaml_url"`
			}{
				Path: ptr("https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package.yaml"),
			},
			wantError: false,
		},
		{
			name: "non-existing URL",
			args: struct {
				Path string `validate:"yaml_url"`
			}{
				Path: "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package",
			},
			wantError: true,
		},
		{
			name: "invalid YAML URL",
			args: struct {
				Path string `validate:"yaml_url"`
			}{
				Path: "raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package.yaml",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			validate := validator.New()
			err := validators.AddYamlUrlValidator(validate)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			err = validate.StructCtx(ctx, tt.args)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
