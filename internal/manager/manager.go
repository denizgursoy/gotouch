//go:generate mockgen -source=./manager.go -destination=mockManager.go -package=manager

package manager

import (
	"io"
	"sync"
)

var (
	manager Manager
	once    = sync.Once{}
)

func GetInstance() Manager {
	once.Do(func() {
		manager = newFileManager()
	})
	return manager
}

type (
	Manager interface {
		CreateDirectoryIfNotExists(directoryName string) error
		GetStream() io.ReadCloser
		IsTest() bool
		GetExtractLocation() string
		GetWd() string
		EditGoModule(projectName, folderName string) error
		AddDependency(dependency string) error
	}
)
