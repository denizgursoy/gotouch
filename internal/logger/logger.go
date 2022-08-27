package logger

import "fmt"

type (
	Logger interface {
		LogErrorIfExists(err error)
		LogInfo(msg string)
	}
	logger struct {
	}
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) LogErrorIfExists(err error) {
	if err != nil {
		fmt.Println(Fata(err.Error()))
	}
}

func (l *logger) LogInfo(msg string) {
	fmt.Println(Info(msg))
}
