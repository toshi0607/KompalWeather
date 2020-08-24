package notifier

import (
	"context"
	"fmt"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/message"
)

// Config

type Slack struct {
}

func (s Slack) Notify(ctx context.Context, result *analyzer.Result) error {
	m := message.Build(result)
	fmt.Print(m)
	return nil
}
