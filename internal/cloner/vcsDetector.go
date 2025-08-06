//go:generate mockgen -source=$GOFILE -destination=mockVCSDetector.go -package=cloner --typed

package cloner

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type VCS string

const (
	VCSNone VCS = ""
	VCSGit  VCS = "git"
)

const (
	gitUploadPackContentType   = "application/x-git-upload-pack-advertisement"
	gitDiscoveryReferencesPath = "/info/refs?service=git-upload-pack"
)

type VCSDetector interface {
	DetectVCS(ctx context.Context, path string) (VCS, error)
}

type RequestExecutor interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ VCSDetector = (*defaultVCSDetector)(nil)

func NewDefaultVCSDetector() defaultVCSDetector {
	return defaultVCSDetector{
		requester: auth.NewAuthenticatedHTTPClient(),
	}
}

type defaultVCSDetector struct {
	requester RequestExecutor
}

// DetectVCS implements VCSDetector.
func (d defaultVCSDetector) DetectVCS(ctx context.Context, rawURL string) (VCS, error) {
	// Using the default remote options
	if IsSSH(rawURL) {
		return checkByListingRemotes(rawURL)
	}

	return d.checkByHttpCall(ctx, rawURL)
}

func checkByListingRemotes(rawURL string) (VCS, error) {
	remoteName := "origin"
	if changedRemote := strings.TrimSpace(os.Getenv("REMOTE_NAME")); len(changedRemote) > 0 {
		remoteName = changedRemote
	}
	remotes := git.NewRemote(nil, &config.RemoteConfig{
		Name: remoteName,
		URLs: []string{rawURL},
	})
	refs, err := remotes.List(&git.ListOptions{
		Auth: getAuth(rawURL), // Uncomment if needed
	})
	if err != nil {
		return VCSNone, err
	}

	// If references are returned, repository exists
	if len(refs) > 0 {
		return VCSGit, nil
	}

	return VCSNone, errors.New("could not find any git repository")
}

func (d defaultVCSDetector) checkByHttpCall(ctx context.Context, rawURL string) (VCS, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", rawURL, gitDiscoveryReferencesPath),
		http.NoBody,
	)
	if err != nil {
		return VCSNone, err
	}

	resp, err := d.requester.Do(req)
	if err != nil {
		return VCSNone, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return VCSNone, nil
	}

	if resp.Header.Get("Content-Type") != gitUploadPackContentType {
		return VCSNone, nil
	}

	return VCSGit, nil
}
