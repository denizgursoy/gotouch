//go:generate mockgen -source=./commander.go -destination=mockCommander.go -package=commander

package commander

import (
	"sync"
)

type (
	Commander interface {
		CreateNewProject(opts *CreateCommandOptions) error
		CompressDirectory(opts *PackageCommandOptions) error
	}
	cmdExecutor struct {
	}
)

var (
	exec Commander
	once sync.Once
)

func GetInstance() Commander {
	once.Do(func() {
		exec = &cmdExecutor{}
	})
	return exec
}
