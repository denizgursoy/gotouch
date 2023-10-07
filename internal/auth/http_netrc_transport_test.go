package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/jdx/go-netrc"

	"github.com/denizgursoy/gotouch/internal/auth"
)

const (
	sampleUsername = "some_user"
	samplePassword = "ar0isheJ"
)

func TestNetrcTransport_RoundTrip_DummyNetrcFile(t *testing.T) {
	testHome := t.TempDir()
	t.Setenv("HOME", testHome)

	testServer := httptest.NewServer(validateBasicAuthHandler(t))
	serverUrl, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("failed to parse server url: %v", err)
	}

	writeNetrcFile(t, filepath.Join(testHome, ".netrc"), serverUrl.Host)
	httpClient := auth.NewAuthenticatedHTTPClient()
	if err != nil {
		t.Fatalf("failed to create authenticated http client: %v", err)
	}

	req, err := http.NewRequest("GET", testServer.URL+"/hello_world", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestNetrcTransport_RoundTrip_NoNetRCFile(t *testing.T) {
	testHome := t.TempDir()
	t.Setenv("HOME", testHome)

	testServer := httptest.NewServer(noOpHttpHandler())
	httpClient := auth.NewAuthenticatedHTTPClient()

	req, err := http.NewRequest("GET", testServer.URL+"/hello_world", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func writeNetrcFile(tb testing.TB, filePath, host string) {
	tb.Helper()
	netRc := netrc.New(filePath)
	netRc.AddMachine(host, sampleUsername, samplePassword)
	if err := netRc.Save(); err != nil {
		tb.Fatalf("failed to save netrc file: %v", err)
	}
}

func noOpHttpHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
}

func validateBasicAuthHandler(tb testing.TB) http.Handler {
	tb.Helper()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok {
			tb.Error("no basic auth set")
		}

		if username != sampleUsername || password != samplePassword {
			tb.Errorf("expected username %s and password %s but got %s and %s", sampleUsername, samplePassword, username, password)
		}

		writer.WriteHeader(http.StatusOK)
	})
}
