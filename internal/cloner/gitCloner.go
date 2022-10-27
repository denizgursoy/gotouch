package cloner

import (
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
)

const (
	GitDirectory = ".git"
)

type (
	gitCloner struct {
		Store store.Store `validate:"required"`
	}
)

func newCloner() Cloner {
	return &gitCloner{
		Store: store.GetInstance(),
	}
}

func (g *gitCloner) CloneFromUrl(url string) error {
	projectName := g.Store.GetValue(store.ProjectName)

	_, err := git.PlainClone(projectName, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	gitDirectory := projectName + string(filepath.Separator) + GitDirectory
	if err = os.RemoveAll(gitDirectory); err != nil {
		return err
	}

	return err
}
