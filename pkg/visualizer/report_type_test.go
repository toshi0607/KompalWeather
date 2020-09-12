package visualizer

import "testing"

func TestReportType_String(t *testing.T) {
	got := DailyReport.String()
	want := "daily"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}
