// go:build !windows
//go:build !windows

package requirements

import (
	"github.com/denizgursoy/gotouch/internal/commandrunner"
)

func getCommand() *commandrunner.CommandData {
	return &commandrunner.CommandData{
		Command: "sh",
		Args:    []string{"init.sh"},
	}
}
