package time

import (
	"fmt"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// https://golang.org/pkg/time/#Time.String
// 保存時に Asia/TokyoになるがParseできない
//var layout = "2006-01-02 15:04:05.999999999 -0700 MST"

var layout = time.RFC3339

// ToJST converts time.Time to JST time
func ToJST(t time.Time) time.Time {
	return t.In(jst)
}

// ToJSTString converts time.Time to string
func ToJSTString(t time.Time) string {
	return t.In(jst).Format(layout)
}

// ToJSTTime converts string to time.Time in JST location
func ToJSTTime(s string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, s, jst)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse in location: %v", err)
	}
	return t, nil
}

type Period struct {
	Start time.Time
	End   time.Time
}

// String returns string expression of Period
func (p Period) String() string {
	format := "2006-01-02"
	prevDate := p.End.AddDate(0, 0, -1)
	return fmt.Sprintf("%s-%s", p.Start.Format(format), prevDate.Format(format))
}

func TodayPeriod(now time.Time) Period {
	tomorrow := now.AddDate(0, 0, 1)
	return Period{
		Start: date(now.Year(), now.Month(), now.Day()),
		End:   date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day()),
	}
}

func YesterdayPeriod(now time.Time) Period {
	yesterday := now.AddDate(0, 0, -1)
	return Period{
		Start: date(yesterday.Year(), yesterday.Month(), yesterday.Day()),
		End:   date(now.Year(), now.Month(), now.Day()),
	}
}

func WeeklyPeriod(now time.Time) Period {
	aWeekAgo := now.AddDate(0, 0, -7)
	return Period{
		Start: date(aWeekAgo.Year(), aWeekAgo.Month(), aWeekAgo.Day()),
		End:   date(now.Year(), now.Month(), now.Day()),
	}
}

func MonthlyPeriod(now time.Time) Period {
	yesterday := now.AddDate(0, 0, -1)
	return Period{
		Start: date(yesterday.Year(), yesterday.Month(), 1),
		End:   date(now.Year(), now.Month(), now.Day()),
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, jst)
}

func NowJST() time.Time {
	return time.Now().In(jst)
}
