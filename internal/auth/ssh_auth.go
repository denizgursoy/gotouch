package auth

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const GitUser = "git"

var defaultKeyFiles = []string{
	"id_rsa",
	"id_ecdsa",
	"id_ecdsa_sk",
	"id_ed25519",
	"id_ed25519_sk",
	"id_dsa",
}

// FindFirstAvailableSSHKey returns an AuthMethod if a default SSH key is found, otherwise nil
func FindFirstAvailableSSHKey() ssh.AuthMethod {
	// Check environment variable first
	password := os.Getenv("SSH_PASSWORD")
	sshKeyFilePath := strings.TrimSpace(os.Getenv("GIT_SSH_KEY"))
	if sshKeyFilePath != "" {
		sshAuth := getSshAuthFromFile(sshKeyFilePath, password)
		if sshAuth != nil {
			return sshAuth
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	for _, keyName := range defaultKeyFiles {
		sshKeyFilePath := filepath.Join(sshDir, keyName)
		sshAuth := getSshAuthFromFile(sshKeyFilePath, password)
		if sshAuth != nil {
			return sshAuth
		}

		continue
	}

	// No key found, return nil (so no SSH auth will be used)
	return nil
}

func getSshAuthFromFile(keyFilePath, password string) ssh.AuthMethod {
	if _, err := os.Stat(keyFilePath); err == nil {
		if auth, err := ssh.NewPublicKeysFromFile(GitUser, keyFilePath, password); err == nil {
			auth.HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
				HostKeyCallback: gossh.InsecureIgnoreHostKey(),
			}

			return auth
		}
	}

	return nil
}
