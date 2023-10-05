package auth

import (
	"os"
	"path/filepath"

	"github.com/jdx/go-netrc"
)

func NetrcFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, netrcFileName), nil
}

func ParseNetrc() (*netrc.Netrc, error) {
	filePath, err := NetrcFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); err != nil {
		return netrc.New(filePath), nil
	}

	return netrc.Parse(filePath)
}
