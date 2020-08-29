package logger

import (
	"fmt"
	"log"
	"os"
)

var _ Logger = (*Log)(nil)

type Log struct {
	logger *log.Logger
}

func NewLog() *Log {
	return &Log{logger: log.New(os.Stdout, "", log.Ltime|log.Llongfile)}
}

func (l Log) SetHandlerName(name string) {
	l.logger.SetPrefix(fmt.Sprintf("[%s] ", name))
}

func (l Log) Info(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l Log) Error(msg string, err error) {
	l.logger.Printf("%s%s", msg, err)
}

func (l Log) Close() error {
	l.logger.Print("Log logger closed")
	return nil
}
