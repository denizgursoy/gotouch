package auth_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/denizgursoy/gotouch/internal/auth"
)

func TestGitNetrcHTTPAuth_SetAuth_DummyNetrcFile(t *testing.T) {
	tests := []struct {
		name           string
		writeNetrcFile bool
		wantAuth       bool
		host           string
		url            string
	}{
		{
			name:           "Matching entry",
			writeNetrcFile: true,
			wantAuth:       true,
			host:           "github.com",
			url:            "https://github.com/denizgursoy/gotouch.git",
		},
		{
			name:           "No matching entry",
			writeNetrcFile: true,
			wantAuth:       false,
			host:           "gitlab.com",
			url:            "https://github.com/denizgursoy/gotouch.git",
		},
		{
			name:           "No netrc file",
			writeNetrcFile: false,
			wantAuth:       false,
			host:           "github.com",
			url:            "https://github.com/denizgursoy/gotouch.git",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			testHome := t.TempDir()
			t.Setenv("HOME", testHome)

			if tt.writeNetrcFile {
				writeNetrcFile(t, filepath.Join(testHome, ".netrc"), tt.host)
			}

			authMethod := auth.NewGitNetrcHTTPAuth()

			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			authMethod.SetAuth(req)

			user, pw, ok := req.BasicAuth()
			if ok != tt.wantAuth {
				t.Fatalf("Expected authentication to be %t but got %t", tt.wantAuth, ok)
			}

			if tt.wantAuth && (user != sampleUsername || pw != samplePassword) {
				t.Errorf("Expected username %s and password %s, but got %s and %s", sampleUsername, samplePassword, user, pw)
			}
		})
	}

}
