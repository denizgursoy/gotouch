package template

import (
	"bytes"
	"testing"
)

func TestSet(t *testing.T) {
	type args struct {
		NoSprig    bool
		LeftDelim  string
		RightDelim string
	}

	tests := []struct {
		name    string
		args    args
		values  any
		content []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "default",
			values: map[string]any{
				"Name": "Deniz",
			},
			content: []byte("Hello {{ .Name }}"),
			want:    []byte("Hello Deniz"),
			wantErr: false,
		},
		{
			name: "delimiter change",
			args: args{
				NoSprig:    false,
				LeftDelim:  "<<",
				RightDelim: ">>",
			},
			values:  map[string]any{"Name": "Deniz"},
			content: []byte("Hello << .Name >>"),
			want:    []byte("Hello Deniz"),
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				NoSprig:    true,
				LeftDelim:  "{{",
				RightDelim: "}}",
			},
			values:  "Deniz",
			content: []byte("Hello {{ . }}"),
			want:    []byte("Hello Deniz"),
			wantErr: false,
		},
		{
			name: "base64",
			args: args{
				NoSprig:    false,
				LeftDelim:  "{{",
				RightDelim: "",
			},
			values:  "RGVuaXo=",
			content: []byte("Hello {{ b64dec . }}"),
			want:    []byte("Hello Deniz"),
			wantErr: false,
		},
	}

	template := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.args.NoSprig {
				template.SetSprigFuncs()
			}
			template.SetDelims(tt.args.LeftDelim, tt.args.RightDelim)

			if v, err := template.Execute(tt.values, string(tt.content)); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			} else if !bytes.Equal([]byte(v), tt.want) {
				t.Errorf("Execute() = %v, want %v", v, tt.want)
			}

			var b bytes.Buffer
			err := template.ExecuteContent(&b, tt.values, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(b.Bytes(), tt.want) {
				t.Errorf("Execute() got = %v, want %v", b.String(), string(tt.want))
			}
		})
	}
}
