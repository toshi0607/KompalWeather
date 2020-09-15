package visualizer

import (
	"errors"

	pt "github.com/toshi0607/kompal-weather/internal/time"
)

type ReportType string

const (
	DailyReport   = ReportType("daily")
	WeeklyReport  = ReportType("weekly")
	MonthlyReport = ReportType("monthly")
	WeekAgoReport = ReportType("weekAgo")
)

func (rt ReportType) String() string {
	return string(rt)
}

func (rt ReportType) IsValid() bool {
	return rt == DailyReport || rt == WeeklyReport || rt == MonthlyReport || rt == WeekAgoReport
}

func (rt ReportType) reportPeriod() (pt.Period, error) {
	switch rt {
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
