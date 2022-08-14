package req

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"reflect"
	"testing"
)

func TestProjectNameRequirement_AskForInput(t *testing.T) {
	tests := []struct {
		name    string
		want    model.Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProjectNameRequirement{}
			got, err := p.AskForInput()
			if (err != nil) != tt.wantErr {
				t.Errorf("AskForInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AskForInput() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_projectNameTask_Complete(t *testing.T) {
	type fields struct {
		ProjectName string
	}
	type args struct {
		in0 interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := projectNameTask{
				ProjectName: tt.fields.ProjectName,
			}
			if got, _ := p.Complete(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Complete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateProjectName(t *testing.T) {
	type args struct {
		projectName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success test 1",
			args:    args{projectName: "github.com/test/project"},
			wantErr: false,
		},
		{
			name:    "success test 2",
			args:    args{projectName: "github.com/test.com/project"},
			wantErr: false,
		},
		{
			name:    "success test 3",
			args:    args{projectName: "github123.com/test123.com/project123"},
			wantErr: false,
		},
		{
			name:    "success test 4",
			args:    args{projectName: "github123.com/test123.com/project123/project"},
			wantErr: false,
		},
		{
			name:    "success test 5",
			args:    args{projectName: "github123"},
			wantErr: false,
		},
		{
			name:    "success test 6",
			args:    args{projectName: "github"},
			wantErr: false,
		},
		{
			name:    "error test 1",
			args:    args{projectName: ""},
			wantErr: true,
		},
		{
			name:    "error test 2",
			args:    args{projectName: "."},
			wantErr: true,
		},
		{
			name:    "error test 3",
			args:    args{projectName: ".exe"},
			wantErr: true,
		},
		{
			name:    "error test 4",
			args:    args{projectName: "./test"},
			wantErr: true,
		},
		{
			name:    "error test 5",
			args:    args{projectName: "123test"},
			wantErr: true,
		},
		{
			name:    "error test 6",
			args:    args{projectName: "error.com/123"},
			wantErr: true,
		},
		{
			name:    "error test 7",
			args:    args{projectName: "error.com/test123/."},
			wantErr: true,
		},
		{
			name:    "error test 8",
			args:    args{projectName: "error.com/test123/blabla."},
			wantErr: true,
		},
		{
			name:    "error test 9",
			args:    args{projectName: "error.com/test123/blabla.exe"},
			wantErr: true,
		},
		{
			name:    "error test 10",
			args:    args{projectName: "error.com/test123.com"},
			wantErr: true,
		},
		{
			name:    "error test 11",
			args:    args{projectName: "error.111/test"},
			wantErr: true,
		},
		{
			name:    "error test 12",
			args:    args{projectName: "error.111"},
			wantErr: true,
		},
		{
			name:    "error test 13",
			args:    args{projectName: "error/errr./test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateProjectName(tt.args.projectName); (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
