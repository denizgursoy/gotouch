package cloner

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
	authMethod, err := g.getAuthMethod(rawUrl)
	if err != nil {
		return err
	}
	cloneOptions := &git.CloneOptions{
		Depth:    1,
		URL:      rawUrl,
		Progress: os.Stdout,
		Auth:     authMethod,
	}

	if len(strings.TrimSpace(branchName)) != 0 {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branchName)
		cloneOptions.SingleBranch = true
		g.Logger.LogInfo(fmt.Sprintf("Cloning branch %s from   -> %s", branchName, rawUrl))
	} else {
		g.Logger.LogInfo("Cloning repository  -> " + rawUrl)
	}

	_, err = git.PlainClone(projectFullPath, false, cloneOptions)
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

func (g *gitCloner) getAuthMethod(rawUrl string) (transport.AuthMethod, error) {
	if isSSHURL(rawUrl) {
		return getSshAuth()
	} else {
		return getNetrcAuth(rawUrl)
	}
}

func isSSHURL(url string) bool {
	return len(url) > 4 && url[:4] == "git@"
}

func getSshAuth() (transport.AuthMethod, error) {
	path, err := getDefaultPrivateKeyPath()
	if err != nil {
		return nil, err
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", path, "")
	if err != nil {
		return nil, err
	}

	return publicKeys, nil
}

func getNetrcAuth(rawUrl string) (transport.AuthMethod, error) {
	gitURL, urlParseError := url.Parse(rawUrl)
	if urlParseError != nil {
		return nil, urlParseError
	}
	switch gitURL.Scheme {
	case "http", "https":
		return auth.NewGitNetrcHTTPAuth(), nil
	}

	return nil, nil
}

func getDefaultPrivateKeyPath() (string, error) {
	// TODO add environment variable
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine user home directory: %v", err)
	}

	keyPath := filepath.Join(homeDir, ".ssh", "id_rsa")

	// Check if the file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", fmt.Errorf("private key not found at %s", keyPath)
	}

	return keyPath, nil
}
