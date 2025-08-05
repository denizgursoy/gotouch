//go:generate mockgen -source=$GOFILE -destination=mockVCSDetector.go -package=cloner --typed

package cloner

import (
	"context"
	"fmt"
	"net/http"

	"github.com/denizgursoy/gotouch/internal/auth"
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
	isGit, err := d.isGit(ctx, rawURL)
	if err != nil {
		return VCSNone, err
	}

	if isGit {
		return VCSGit, nil
	}

	return VCSNone, nil
}

func (d defaultVCSDetector) isGit(ctx context.Context, rawURL string) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", rawURL, gitDiscoveryReferencesPath),
		http.NoBody,
	)
	if err != nil {
		return false, err
	}

	resp, err := d.requester.Do(req)
	if err != nil {
		return false, nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false, nil
	}

	if resp.Header.Get("Content-Type") != gitUploadPackContentType {
		return false, nil
	}

	return true, nil
}
