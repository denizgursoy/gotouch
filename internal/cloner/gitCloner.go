package cloner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

const (
	GitDirectory = ".git"
)

type (
	gitCloner struct {
		Store  store.Store   `validate:"required"`
		Logger logger.Logger `validate:"required"`
	}
)

func newCloner() Cloner {
	return &gitCloner{
		Store:  store.GetInstance(),
		Logger: logger.NewLogger(),
	}
}

func (g *gitCloner) CloneFromUrl(ctx context.Context, rawUrl, branchName string) error {
	projectFullPath := g.Store.GetValue(store.ProjectFullPath)

	cloneOptions := &git.CloneOptions{
		Depth:    1,
		URL:      rawUrl,
		Progress: os.Stdout,
		Auth:     getAuth(rawUrl),
	}

	if len(strings.TrimSpace(branchName)) != 0 {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branchName)
		cloneOptions.SingleBranch = true
		g.Logger.LogInfo(fmt.Sprintf("Cloning branch %s from   -> %s", branchName, rawUrl))
	} else {
		g.Logger.LogInfo("Cloning repository  -> " + rawUrl)
	}

	_, err := git.PlainCloneContext(ctx, projectFullPath, false, cloneOptions)
	if err != nil {
		if !errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return err
		}
	}

	gitDirectory := filepath.Join(projectFullPath, GitDirectory)
	if err = os.RemoveAll(gitDirectory); err != nil {
		return err
	}
	g.Logger.LogInfo("Cloned successfully")

	return err
}

func getAuth(rawUrl string) transport.AuthMethod {
	if IsSSH(rawUrl) {
		return auth.FindFirstAvailableSSHKey()
	} else if IsHTTP(rawUrl) {
		return auth.NewGitNetrcHTTPAuth()
	}

	return nil
}

func IsSSH(url string) bool {
	return strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://")
}

func IsHTTP(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
