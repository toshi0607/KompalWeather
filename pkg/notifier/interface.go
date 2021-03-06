package notifier

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

// Notifier is an interface of notifier
type Notifier interface {
	Notify(ctx context.Context, result *analyzer.Result) error
	NotifyWithMedium(ctx context.Context, msg string, contents [][]byte) error
	Type() string
}
