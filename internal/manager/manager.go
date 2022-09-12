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
		CreateFile(reader io.ReadCloser, path string) error
		EditGoModule() error
		GetExtractLocation() string
	}
)
