package kompal

import "github.com/toshi0607/kompal-weather/pkg/status"

type fetcher interface {
	Fetch() (*status.Status, error)
}
