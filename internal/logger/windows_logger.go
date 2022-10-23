// go:build windows
// +build windows

package logger

import "github.com/sirupsen/logrus"

func logError(msg string) {
	logrus.Error(msg)
}

func logInfo(msg string) {
	logrus.Info(msg)
}
