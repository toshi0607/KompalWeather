package logger

import (
	"fmt"
	"log"
	"os"
)

var _ Logger = (*Log)(nil)

// Log represents log
type Log struct {
	logger *log.Logger
}

// NewLog builds new Log
func NewLog() *Log {
	return &Log{logger: log.New(os.Stdout, "", log.Ltime|log.Llongfile)}
}

// SetHandlerName sets prefix to logger
func (l Log) SetHandlerName(name string) {
	l.logger.SetPrefix(fmt.Sprintf("[%s] ", name))
}

// Info prints information
func (l Log) Info(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

// Error prints error string
func (l Log) Error(msg string, err error) {
	l.logger.Printf("%s%s", msg, err)
}

// Close prints log
func (l Log) Close() error {
	l.logger.Print("Log logger closed")
	return nil
}
