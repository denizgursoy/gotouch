//go:generate mockgen -source=./manager.go -destination=mock-manager.go -package=manager

package manager

import (
	"io"
)

var (
	manager Manager
)

func GetInstance() Manager {
	return manager
}

type (
	Manager interface {
		CreateDirectoryIfNotExists(directoryName string) error
		GetStream() io.ReadCloser
		IsTest() bool
		GetExtractLocation() string
		GetWd() string
	}

	fManager struct {
	}
)
