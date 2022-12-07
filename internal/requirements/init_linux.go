package requirements

import (
	"github.com/denizgursoy/gotouch/internal/store"
)

const InitFileName = "init.sh"

func executeInitFile(str store.Store) error {
	commandData := CommandData{
		Command: "sh",
		Args:    []string{"init.sh"},
	}
	return RunCommand(&commandData, str)
}
