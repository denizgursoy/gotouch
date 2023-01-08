package logger

import "fmt"

type (
	Logger interface {
		LogErrorIfExists(err error)
		LogInfo(msg string)
	}
	logger struct{}
)

var (
	Red  = Color("\033[1;31m%s\033[0m")
	Teal = Color("\033[1;36m%s\033[0m")
)

func Color(colorString string) func(...any) string {
	sprint := func(args ...any) string {
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
		logError(Red(err.Error()))
	}
}

func (l *logger) LogInfo(msg string) {
	logInfo(Teal(msg))
}
