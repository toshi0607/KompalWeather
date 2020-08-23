package analyzer

import (
	"errors"

	"github.com/toshi0607/kompal-weather/pkg/status"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"golang.org/x/net/context"
)

type analyzer struct {
	storage storage.Storage
}

type (
	Result struct {
		maleTrend    Trend
		femaleTrend  Trend
		latestStatus status.Status
	}

	Trend int
)

const (
	Unknown    = 0
	Increasing = 1
	Decreasing = 2
	Constant   = 3
)

func (t Trend) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Increasing:
		return "Increasing"
	case Decreasing:
		return "Decreasing"
	case Constant:
		return "Constant"
	default:
		return "Invalid"
	}
}

const expectedNumberOfStatuses = 2

func New(storage storage.Storage) Analyzer {
	return analyzer{
		storage: storage,
	}
}

func (a analyzer) Analyze(ctx context.Context) (*Result, error) {
	ss, err := a.storage.Statuses(ctx)
	if err != nil {
		return nil, err
	}
	if len(ss) != expectedNumberOfStatuses {
		return nil, errors.New("no sufficient status")
	}

	// Expect time series to be ss[0],ss[1]
	if ss[0].CreatedAt.After(ss[1].CreatedAt) {
		ss[0], ss[1] = ss[1], ss[0]
	}

	var result Result
	if ss[0].MaleSauna == ss[1].MaleSauna {
		result.maleTrend = Constant
	} else if ss[0].MaleSauna > ss[1].MaleSauna {
		result.maleTrend = Decreasing
	} else {
		result.maleTrend = Increasing
	}

	if ss[0].FemaleSauna == ss[1].FemaleSauna {
		result.femaleTrend = Constant
	} else if ss[0].FemaleSauna > ss[1].FemaleSauna {
		result.femaleTrend = Decreasing
	} else {
		result.femaleTrend = Increasing
	}

	result.latestStatus = ss[1]

	return &result, nil
}
