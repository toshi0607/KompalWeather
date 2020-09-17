package path

import (
	"fmt"
	"strings"
	"testing"

	pt "github.com/toshi0607/kompal-weather/internal/time"
	"github.com/toshi0607/kompal-weather/pkg/report"
)

var (
	format    = "2006-01-02"
	now       = pt.NowJST().Format(format)
	weekAgo   = pt.NowJST().AddDate(0, 0, -7).Format(format)
	yesterday = pt.NowJST().AddDate(0, 0, -1).Format(format)
)

func TestReportObject(t *testing.T) {
	tests := map[string]struct {
		k      report.Kind
		gender string
		want   string
	}{
		"daily": {
			k:      report.DailyReport,
			gender: "male",
			want:   fmt.Sprintf("daily/%s-%s-male.png", now, now),
		},
		"weekAgo": {
			k:      report.WeekAgoReport,
			gender: "female",
			want:   fmt.Sprintf("daily/%s-%s-female.png", weekAgo, weekAgo),
		},
	}

	for name, te := range tests {
		te := te

		got, err := ReportObject(te.gender, te.k)
		if err != nil {
			t.Error(err)
		}

		if got != te.want {
			t.Errorf("[%s] got: %s, want: %s", name, got, te.want)
		}
	}
}

func TestMaleWeekAgoReportObject(t *testing.T) {
	want := fmt.Sprintf("daily/%s-%s-male.png", weekAgo, weekAgo)
	got := MaleWeekAgoReportObject()
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestFemaleWeekAgoReportObject(t *testing.T) {
	want := fmt.Sprintf("daily/%s-%s-female.png", weekAgo, weekAgo)
	got := FemaleWeekAgoReportObject()
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestMaleWeeklyReportObject(t *testing.T) {
	want := fmt.Sprintf("weekly/%s-%s-male.png", weekAgo, yesterday)
	got := MaleWeeklyReportObject()
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestFemaleWeeklyReportObject(t *testing.T) {
	want := fmt.Sprintf("weekly/%s-%s-female.png", weekAgo, yesterday)
	got := FemaleWeeklyReportObject()
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestLogObject(t *testing.T) {
	fileName := "test.png"
	wantPrefix := "logs/"
	wantSuffix := fileName

	got := LogObject(fileName)

	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("got: %s, wantPrefix: %s", got, wantPrefix)
	}
	if !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("got: %s, wantSuffix: %s", got, wantSuffix)
	}
}
