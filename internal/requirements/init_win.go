// go:build windows
//go:build windows
// +build windows

package requirements

import "github.com/denizgursoy/gotouch/internal/commandrunner"

func getCommand() *commandrunner.CommandData {
	return &commandrunner.CommandData{
		Command: "CMD",
		Args: []string{
			"/C",
			WindowsInitFile,
		},
	}
}
