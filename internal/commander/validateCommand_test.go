package commander

import (
	"testing"

	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/golang/mock/gomock"
)

func Test_cmdExecutor_ValidateYaml(t *testing.T) {
	type args struct {
		opts *ValidateCommandOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmdExecutor{}
			if err := c.ValidateYaml(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("ValidateYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValidYaml(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockLister := lister.NewMockLister(controller)
	mockLogger := logger.NewLogger()
	path := "../testdata/input.yaml"
	failPath := "../testdata/xxx.yaml"

	type args struct {
		opts *ValidateCommandOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully",
			args: args{
				opts: &ValidateCommandOptions{
					Lister: mockLister,
					Logger: mockLogger,
					Path:   &path,
				},
			},
			wantErr: false,
		},
		{
			name: "missing invalid path",
			args: args{
				opts: &ValidateCommandOptions{
					Lister: mockLister,
					Logger: mockLogger,
					Path:   &failPath,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValidYaml(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("isValidYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
