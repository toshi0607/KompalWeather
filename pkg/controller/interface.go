package controller

import (
	"context"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/report"
)

// Controller is an interface of a controller
type Controller interface {
	Watch(ctx context.Context) (*analyzer.Result, error)
	Trend(ctx context.Context, k report.Kind) error
}
