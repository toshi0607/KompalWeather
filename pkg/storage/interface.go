package storage

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/status"
)

// Storage is an interface of storage
type Storage interface {
	Statuses(ctx context.Context) ([]status.Status, error)
	Save(ctx context.Context, st *status.Status) (*status.Status, error)
}
