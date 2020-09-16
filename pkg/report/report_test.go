package report

import "testing"

func TestKind_String(t *testing.T) {
	got := DailyReport.String()
	want := "daily"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}

	tests := map[string]struct {
		k    Kind
		want string
	}{
		"valid daily":   {k: DailyReport, want: "daily"},
		"valid weekly":  {k: WeeklyReport, want: "weekly"},
		"valid monthly": {k: MonthlyReport, want: "monthly"},
		"valid weekAgo": {k: WeekAgoReport, want: "weekAgo"},
		"invalid":       {k: Kind("invalid type"), want: "invalid type"},
	}

	for name, te := range tests {
		te := te
		got := te.k.String()
		want := te.want
		if got != want {
			t.Errorf("[%s] got: %v, want: %v", name, got, want)
		}
	}
}

func TestKind_IsValid(t *testing.T) {
	tests := map[string]struct {
		k    Kind
		want bool
	}{
		"valid daily":   {k: DailyReport, want: true},
		"valid weekly":  {k: WeeklyReport, want: true},
		"valid monthly": {k: MonthlyReport, want: true},
		"valid weekAgo": {k: WeekAgoReport, want: true},
		"invalid":       {k: Kind("invalid type"), want: false},
	}

	for name, te := range tests {
		te := te
		got := te.k.IsValid()
		want := te.want
		if got != want {
			t.Errorf("[%s] got: %v, want: %v", name, got, want)
		}
	}
}
