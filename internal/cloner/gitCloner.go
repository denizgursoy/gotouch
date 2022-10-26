package cloner

import (
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-git/go-git/v5"
	"os"
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

	return err
}
