package visualizer

import (
	"testing"
)

func TestReportType_String(t *testing.T) {
	got := DailyReport.String()
	want := "daily"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}

	tests := map[string]struct {
		rt   ReportType
		want string
	}{
		"valid daily":   {rt: DailyReport, want: "daily"},
		"valid weekly":  {rt: WeeklyReport, want: "weekly"},
		"valid monthly": {rt: MonthlyReport, want: "monthly"},
		"valid weekAgo": {rt: WeekAgoReport, want: "weekAgo"},
		"invalid":       {rt: ReportType("invalid type"), want: "invalid type"},
	}

	for name, te := range tests {
		te := te
		got := te.rt.String()
		want := te.want
		if got != want {
			t.Errorf("[%s] got: %v, want: %v", name, got, want)
		}
	}
}

func TestReportType_IsValid(t *testing.T) {
	tests := map[string]struct {
		rt   ReportType
		want bool
	}{
		"valid daily":   {rt: DailyReport, want: true},
		"valid weekly":  {rt: WeeklyReport, want: true},
		"valid monthly": {rt: MonthlyReport, want: true},
		"valid weekAgo": {rt: WeekAgoReport, want: true},
		"invalid":       {rt: ReportType("invalid type"), want: false},
	}

	for name, te := range tests {
		te := te
		got := te.rt.IsValid()
		want := te.want
		if got != want {
			t.Errorf("[%s] got: %v, want: %v", name, got, want)
		}
	}
}
