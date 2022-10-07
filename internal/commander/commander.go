//go:generate mockgen -source=./commander.go -destination=mockCommander.go -package=commander

package commander

import (
	"sync"
)

type (
	Commander interface {
		CreateNewProject(opts *CreateCommandOptions) error
		CompressDirectory(opts *PackageCommandOptions) error
		ValidateYaml(opts *ValidateCommandOptions) error
	}
	cmdExecutor struct {
	}
)

var (
	exc  Commander
	once sync.Once
)

func GetInstance() Commander {
	once.Do(func() {
		exc = &cmdExecutor{}
	})
	return exc
}
