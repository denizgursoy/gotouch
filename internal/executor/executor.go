//go:generate mockgen -source=./executor.go -destination=mockExecutor.go -package=executor

package executor

import (
	"context"
	"sync"

	"github.com/denizgursoy/gotouch/internal/model"
)

type (
	Tasks []model.Task

	Requirements []model.Requirement

	Executor interface {
		Execute(ctx context.Context, requirements Requirements) error
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
