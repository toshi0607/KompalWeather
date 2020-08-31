package kompal

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

// Fetcher is an interface of a fetcher
type Fetcher interface {
	Fetch(ctx context.Context) (*status.Status, error)
}
