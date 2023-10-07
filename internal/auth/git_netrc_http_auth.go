package auth

import (
	"fmt"
	"net/http"

	httptransport "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jdx/go-netrc"
)

var _ httptransport.AuthMethod = (*GitNetrcHTTPAuth)(nil)

const (
	netrcLoginKey    = "login"
	netrcPasswordKey = "password"
	netrcAuthName    = "http-netrc-auth"
	netrcEmptyString = "netrc://<empty>"
)

func NewGitNetrcHTTPAuth() *GitNetrcHTTPAuth {
	creds, err := ParseNetrc()
	if err != nil {
		return new(GitNetrcHTTPAuth)
	}

	return &GitNetrcHTTPAuth{
		credentials: creds,
	}
}

type GitNetrcHTTPAuth struct {
	credentials *netrc.Netrc
}

func (g GitNetrcHTTPAuth) String() string {
	if g.credentials == nil {
		return netrcEmptyString
	}

	return fmt.Sprintf("netrc://%s", g.credentials.Path)
}

func (GitNetrcHTTPAuth) Name() string {
	return netrcAuthName
}

func (g GitNetrcHTTPAuth) SetAuth(r *http.Request) {
	if g.credentials == nil {
		return
	}

	machine := g.credentials.Machine(r.Host)
	if machine == nil {
		return
	}

	r.SetBasicAuth(machine.Get(netrcLoginKey), machine.Get(netrcPasswordKey))
}
