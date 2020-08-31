package logger

// Logger is an interface of logger
type Logger interface {
	Info(format string, args ...interface{})
	Error(msg string, err error)
	SetHandlerName(name string)
	Close() error
}
