package cloner

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
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

func (g *gitCloner) CloneFromUrl(rawUrl, branchName string) error {
	projectFullPath := g.Store.GetValue(store.ProjectFullPath)

	cloneOptions := &git.CloneOptions{
		Depth:    1,
		URL:      rawUrl,
		Progress: os.Stdout,
	}

	gitURL, urlParseError := url.Parse(rawUrl)
	if urlParseError != nil {
		return urlParseError
	}

	switch gitURL.Scheme {
	case "http", "https":
		cloneOptions.Auth = auth.NewGitNetrcHTTPAuth()
	}

	if len(strings.TrimSpace(branchName)) != 0 {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branchName)
		cloneOptions.SingleBranch = true
		g.Logger.LogInfo(fmt.Sprintf("Cloning branch %s from   -> %s", branchName, rawUrl))
	} else {
		g.Logger.LogInfo("Cloning repository  -> " + rawUrl)
	}

	_, err := git.PlainClone(projectFullPath, false, cloneOptions)
	if err != nil {
		return err
	}

	gitDirectory := filepath.Join(projectFullPath, GitDirectory)
	if err = os.RemoveAll(gitDirectory); err != nil {
		return err
	}
	g.Logger.LogInfo("Cloned successfully")
	return err
}
