package report

import (
	"errors"

	pt "github.com/toshi0607/kompal-weather/internal/time"
)

type Kind string

const (
	DailyReport   = Kind("daily")
	WeeklyReport  = Kind("weekly")
	MonthlyReport = Kind("monthly")
	WeekAgoReport = Kind("weekAgo")
)

func (k Kind) String() string {
	return string(k)
}

func (k Kind) IsValid() bool {
	return k == DailyReport || k == WeeklyReport || k == MonthlyReport || k == WeekAgoReport
}

func (k Kind) Period() (pt.Period, error) {
	switch k {
	case DailyReport:
		return pt.TodayPeriod(pt.NowJST()), nil
	case WeeklyReport:
		return pt.WeeklyPeriod(pt.NowJST()), nil
	case MonthlyReport:
		return pt.MonthlyPeriod(pt.NowJST()), nil
	case WeekAgoReport:
		return pt.WeekAgoPeriod(pt.NowJST()), nil
	default:
		return pt.Period{}, errors.New("unknown report type")
	}
}
