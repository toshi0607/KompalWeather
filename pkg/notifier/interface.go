package notifier

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

type Notifier interface {
	Notify(ctx context.Context, result *analyzer.Result) error
	Type() string
}
