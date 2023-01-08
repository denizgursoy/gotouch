//go:build !windows

package logger

import (
	"fmt"
)

func logError(msg string) {
	fmt.Println(msg)
}

func logInfo(msg string) {
	fmt.Println(msg)
}
