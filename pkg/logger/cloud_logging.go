package logger

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
)

var _ Logger = (*CloudLogging)(nil)

// CloudLogging represents Cloud Logging
type CloudLogging struct {
	client      *logging.Client
	serviceName string
	handlerName string
	version     string
	environment string
}

// Message represents a payload of a logging.Entry
type Message struct {
	Err         string `json:",omitempty"`
	HandlerName string `json:"handlerName,omitempty"`
	Msg         string `json:"msg"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// NewCloudLogging builds new CloudLogging
func NewCloudLogging(ctx context.Context, gcpProjectID, serviceName, version, environment string) (*CloudLogging, error) {
	client, err := logging.NewClient(ctx, gcpProjectID)
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

// Close closes a cloud logging client
func (l *CloudLogging) Close() error {
	return l.client.Close()
}

// SetHandlerName sets a handler name
func (l *CloudLogging) SetHandlerName(name string) {
	l.handlerName = name
}

// Info sends information Cloud Logging
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
