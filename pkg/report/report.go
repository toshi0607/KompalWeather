package report

import (
	"errors"

	pt "github.com/toshi0607/kompal-weather/internal/time"
)

type Kind string

// To handle weekAgo kind transparently, keep period as fields in Report struct
// New creates default and NewWithDate (or option) handles start date

const (
	DailyReport   = Kind("daily")
	WeeklyReport  = Kind("weekly")
	MonthlyReport = Kind("monthly")
	// This is actually a daily report start from 7days ago.
	// On the the other hand, it's a good variant if cloud scheduler isn't good at
	// building a dynamic request body.
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
