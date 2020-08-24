package slack

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

func TestNotify(t *testing.T) {
	ctx := context.TODO()
	url := os.Getenv("WEBHOOK_URL")
	s := New(
		&Config{
			WebhookUrl:   url,
			ChannelNames: []string{"dev"},
			UserName:     "kompal-weather",
		})

	err := s.Notify(ctx, &analyzer.Result{
		MaleTrend:   analyzer.Unknown,
		FemaleTrend: analyzer.Constant,
		LatestStatus: status.Status{
			MaleSauna:   status.Normal,
			FemaleSauna: status.Full,
			Timestamp:   time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}
