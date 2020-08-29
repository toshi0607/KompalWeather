package logger

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/logging"
)

type Log struct {
	client      *logging.Client
	serviceName string
	handlerName string
	version     string
	environment string
}

type Message struct {
	Err         string `json:",omitempty"`
	HandlerName string `json:"handlerName",omitempty`
	Msg         string `json:"msg"`
	Version     string `json:"version"`
}

func New(ctx context.Context, gcpProjectId, serviceName, version, environment string) (*Log, error) {
	client, err := logging.NewClient(ctx, gcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to new logger: %v", err)
	}
	return &Log{
		client:      client,
		serviceName: serviceName,
		version:     version,
		environment: environment,
	}, nil
}

func (l *Log) Close() error {
	return l.client.Close()
}

func (l *Log) SetHandlerName(name string) {
	l.handlerName = name
}

// Need roles/logging.logWriter	to write in Cloud Logging
func (l *Log) Info(format string, args ...interface{}) {
	if l.environment == "local" {
		log.Printf(format, args...)
	}
	l.client.Logger(l.serviceName).Log(logging.Entry{
		Severity: logging.Info,
		Payload: Message{
			Msg:         fmt.Sprintf(format, args...),
			Version:     l.version,
			HandlerName: l.handlerName,
		},
	})
}

func (l *Log) Error(msg string, err error) {
	if l.environment == "local" {
		log.Printf("%s%s", msg, err)
	}
	l.client.Logger(l.serviceName).Log(logging.Entry{
		Severity: logging.Info,
		Payload: Message{
			Msg:         msg,
			Err:         err.Error(),
			Version:     l.version,
			HandlerName: l.handlerName,
		},
	})
}