//go:build windows

package logger

import (
	"io"

	"github.com/mattn/go-colorable"
)

var colorableStdout io.Writer

func init() {
	colorableStdout = colorable.NewColorableStdout()
}

func logError(msg string) {
	colorableStdout.Write([]byte(msg))
}

func logInfo(msg string) {
	colorableStdout.Write([]byte(msg))
}
