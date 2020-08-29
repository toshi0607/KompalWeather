package logger

type Logger interface {
	Info(format string, args ...interface{})
	Error(msg string, err error)
	SetHandlerName(name string)
	Close() error
}
