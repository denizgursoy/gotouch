package logger

type (
	Logger interface {
		LogErrorIfExists(err error)
		LogInfo(msg string)
	}
	logger struct{}
)

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) LogErrorIfExists(err error) {
	if err != nil {
		logError(err.Error())
	}
}

func (l *logger) LogInfo(msg string) {
	logInfo(msg)
}
