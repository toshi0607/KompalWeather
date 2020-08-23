package kompal

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

type Fetcher interface {
	Fetch(ctx context.Context) (*status.Status, error)
}
