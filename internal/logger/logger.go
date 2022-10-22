package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

type (
	Logger interface {
		LogErrorIfExists(err error)
		LogInfo(msg string)
	}
	logger struct{}
)

var (
	Info = Teal
	Fata = Red
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

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) LogErrorIfExists(err error) {
	if err != nil {
		if runtime.GOOS == "windows" {
			logrus.Error(err.Error())
		} else {
			fmt.Println(Fata(err.Error()))
		}
	}
}

func (l *logger) LogInfo(msg string) {
	if runtime.GOOS == "windows" {
		logrus.Info(msg)
	} else {
		fmt.Println(Info(msg))
	}
}
