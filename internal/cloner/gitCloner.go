package cloner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func (g *gitCloner) CloneFromUrl(url, branchName string) error {
	projectName := g.Store.GetValue(store.ProjectName)

	var name plumbing.ReferenceName
	if len(strings.TrimSpace(branchName)) != 0 {
		name = plumbing.NewBranchReferenceName(branchName)
		g.Logger.LogInfo(fmt.Sprintf("Cloning branch %s from   -> %s", branchName, url))
	} else {
		g.Logger.LogInfo("Cloning repository  -> " + url)
	}

	cloneOptions := &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: name,
	}

	_, err := git.PlainClone(projectName, false, cloneOptions)
	if err != nil {
		return err
	}

	gitDirectory := projectName + string(filepath.Separator) + GitDirectory
	if err = os.RemoveAll(gitDirectory); err != nil {
		return err
	}
	g.Logger.LogInfo("Cloned successfully")
	return err
}
