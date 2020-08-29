package analyzer

import (
	"errors"
	"fmt"

	"github.com/toshi0607/kompal-weather/pkg/status"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"golang.org/x/net/context"
)

type analyzer struct {
	storage storage.Storage
}

type (
	Result struct {
		MaleTrend    Trend
		FemaleTrend  Trend
		LatestStatus status.Status
	}

	Trend int
)

const (
	Unknown    = 0
	Increasing = 1
	Decreasing = 2
	Constant   = 3
	Open       = 4
	Close      = 5
)

func (t Trend) String() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Increasing:
		return "混んできました。"
	case Decreasing:
		return "空いてきました。"
	case Constant:
		return "変わりありません。"
	case Open:
		return "営業を開始しました。"
	case Close:
		return "営業を終了しました。"
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
		return nil, fmt.Errorf("failed to get statuses: %v", err)
	}
	if len(ss) != expectedNumberOfStatuses {
		return nil, errors.New("no sufficient status")
	}

	// Expect time series to be ss[0],ss[1]
	if ss[0].CreatedAt.After(ss[1].CreatedAt) {
		ss[0], ss[1] = ss[1], ss[0]
	}

	var result Result
	result.LatestStatus = ss[1]

	if ss[0].MaleSauna != status.Off && ss[1].MaleSauna == status.Off {
		result.MaleTrend = Close
	} else if ss[0].MaleSauna == status.Off && ss[1].MaleSauna != status.Off {
		result.MaleTrend = Open
	} else if ss[0].MaleSauna == ss[1].MaleSauna {
		result.MaleTrend = Constant
	} else if ss[0].MaleSauna > ss[1].MaleSauna {
		result.MaleTrend = Decreasing
	} else {
		result.MaleTrend = Increasing
	}

	if ss[0].FemaleSauna != status.Off && ss[1].FemaleSauna == status.Off {
		result.FemaleTrend = Close
	} else if ss[0].FemaleSauna == status.Off && ss[1].FemaleSauna != status.Off {
		result.FemaleTrend = Open
	} else if ss[0].FemaleSauna == ss[1].FemaleSauna {
		result.FemaleTrend = Constant
	} else if ss[0].FemaleSauna > ss[1].FemaleSauna {
		result.FemaleTrend = Decreasing
	} else {
		result.FemaleTrend = Increasing
	}

	return &result, nil
}
