package util

import "testing"

func TestGetBaseName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success test 1",
			args:    args{path: "example.com/test"},
			want:    "test",
			wantErr: false,
		},
		{
			name:    "success test 2",
			args:    args{path: "folder"},
			want:    "folder",
			wantErr: false,
		},
		{
			name:    ".exe",
			args:    args{path: ".exe"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBaseName(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBaseName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
