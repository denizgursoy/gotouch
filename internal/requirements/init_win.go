// go:build windows
//go:build windows
// +build windows

package requirements

import "github.com/denizgursoy/gotouch/internal/store"

func executeInitFile(str store.Store) error {
	commandData := CommandData{
		Command: InitFileName,
	}
	return RunCommand(&commandData, str)
}
