package util

import (
	"path/filepath"
)

func GetBaseName(path string) (string, error) {
	path = filepath.Base(path)
	// TODO: error handling
	return path, nil
}
