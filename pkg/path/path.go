package path

import (
	"fmt"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/report"
)

// ReportObjectPath returns full path for GCS
// Example:
//   daily:   daily/2020-09-09-2020-09-09-male.png
//   weekly:  weekly/2020-12-28-2021-01-03-female.png
//   monthly: monthly/2020-12-01-2020-12-31-male.png
func ReportObjectPath(gender string, k report.Kind) (string, error) {
	periodStr, err := k.Period()
	if err != nil {
		return "", err
	}
	if k == report.WeekAgoReport {
		k = report.DailyReport
	}
	return fmt.Sprintf("%s/%s-%s.png", k, periodStr, gender), nil
}

// LogObjectPath returns full path for GCS
//   logs:    logs/1599983507/last-page.png
func LogObjectPath(fileName string) string {
	return fmt.Sprintf("logs/%v/%s", time.Now().Unix(), fileName)
}
