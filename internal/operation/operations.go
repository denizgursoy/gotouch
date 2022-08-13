package operation

import (
	"github.com/denizgursoy/gotouch/internal/lister"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/denizgursoy/gotouch/internal/uncompressor"
	"sync"
)

type (
	Tasks []model.Task

	Requirements []model.Requirement

	Executor interface {
		Execute(requirements Requirements) error
	}
)

var (
	Prompter     = prompts.GetInstance()
	Extractor    = uncompressor.GetInstance()
	Lister       = lister.GetInstance()
	MainExecutor Executor
	once         sync.Once
)

func init() {
	MainExecutor = GetInstance()
}

func GetInstance() Executor {
	once.Do(func() {
		MainExecutor = newExecutor()
	})
	return MainExecutor
}
