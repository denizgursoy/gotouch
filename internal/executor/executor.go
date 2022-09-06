//go:generate mockgen -source=./executor.go -destination=mockExecutor.go -package=executor

package executor

import (
	"sync"

	"github.com/denizgursoy/gotouch/internal/model"
)

type (
	Tasks []model.Task

	Requirements []model.Requirement

	CommandData struct {
		WorkingDir *string
		Command    string
		Args       []string
	}

	Executor interface {
		Execute(requirements Requirements) error
		RunCommand(data *CommandData) error
	}
)

var (
	main Executor
	once sync.Once
)

func init() {
	main = GetInstance()
}

func GetInstance() Executor {
	once.Do(func() {
		main = newExecutor()
	})
	return main
}
