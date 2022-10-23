// go:build !windows
// +build !windows

package logger

import (
	"fmt"
)

var (
	Red  = Color("\033[1;31m%s\033[0m")
	Teal = Color("\033[1;36m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func logError(msg string) {
	fmt.Println(Red(msg))
}

func logInfo(msg string) {
	fmt.Println(Teal(msg))
}
