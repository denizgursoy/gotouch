package operation

import (
	"github.com/denizgursoy/gotouch/internal/model"
	"sync"
)

type (
	Tasks []model.Task

	Requirements []model.Requirement

	CommandData struct {
		Command string
		Args    []string
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
