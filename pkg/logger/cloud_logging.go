package logger

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
)

var _ Logger = (*CloudLogging)(nil)

type CloudLogging struct {
	client      *logging.Client
	serviceName string
	handlerName string
	version     string
	environment string
}

type Message struct {
	Err         string `json:",omitempty"`
	HandlerName string `json:"handlerName,omitempty"`
	Msg         string `json:"msg"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func NewCloudLogging(ctx context.Context, gcpProjectId, serviceName, version, environment string) (*CloudLogging, error) {
	client, err := logging.NewClient(ctx, gcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to new logger: %v", err)
	}
	return &CloudLogging{
		client:      client,
		serviceName: serviceName,
		version:     version,
		environment: environment,
	}, nil
}

func (l *CloudLogging) Close() error {
	return l.client.Close()
}

func (l *CloudLogging) SetHandlerName(name string) {
	l.handlerName = name
}

// Need roles/logging.logWriter	to write in Cloud Logging
func (l *CloudLogging) Info(format string, args ...interface{}) {
	l.client.Logger(l.serviceName).Log(logging.Entry{
		Severity: logging.Info,
		Payload: Message{
			Msg:         fmt.Sprintf(format, args...),
			Version:     l.version,
			HandlerName: l.handlerName,
			Environment: l.environment,
		},
	})
}

func (l *CloudLogging) Error(msg string, err error) {
	l.client.Logger(l.serviceName).Log(logging.Entry{
		Severity: logging.Info,
		Payload: Message{
			Msg:         msg,
			Err:         err.Error(),
			Version:     l.version,
			HandlerName: l.handlerName,
			Environment: l.environment,
		},
	})
}
