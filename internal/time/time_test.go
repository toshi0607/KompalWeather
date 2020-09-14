package time

import (
	"testing"
	"time"
)

func TestToJST(t *testing.T) {
	time, err := time.Parse(layout, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Error(err)
	}
	got := ToJST(time)
	if got.Location() != jst {
		t.Errorf("got: %s, want: %s", got.Location(), jst)
	}
}

func TestToJSTString(t *testing.T) {
	time, err := time.Parse(layout, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Error(err)
	}
	got := ToJSTString(time)
	want := "2006-01-03T00:04:05+09:00"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestToJSTTime(t *testing.T) {
	got, err := ToJSTTime("2006-01-03T00:04:05+09:00")
	if err != nil {
		t.Error(err)
	}
	if got.Location() != jst {
		t.Errorf("got: %s, want: %s", got.Location(), jst)
	}
}

func testDate(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, jst)
}

func TestTodayPeriod(t *testing.T) {
	got := TodayPeriod(testDate(2020, 9, 7)).String()
	want := "2020-09-07-2020-09-07"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestYesterdayPeriod(t *testing.T) {
	got := YesterdayPeriod(testDate(2020, 9, 7)).String()
	want := "2020-09-06-2020-09-06"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestWeeklyPeriod(t *testing.T) {
	got := WeeklyPeriod(testDate(2020, 9, 7)).String()
	want := "2020-08-31-2020-09-06"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMonthlyPeriod(t *testing.T) {
	got := MonthlyPeriod(testDate(2020, 10, 1)).String()
	want := "2020-09-01-2020-09-30"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestPeriod_String(t *testing.T) {
	got := Period{
		Start: testDate(2020, 9, 1),
		End:   testDate(2020, 9, 2),
	}.String()
	want := "2020-09-01-2020-09-01"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestNowJST(t *testing.T) {
	got := NowJST()
	if got.Location() != jst {
		t.Errorf("got: %s, want: %s", got.Location(), jst)
	}
	if got.IsZero() {
		t.Errorf("Timestamp should not be zero value")
	}
}
