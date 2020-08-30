package slack

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

func TestSlack_Type(t *testing.T) {
	got := New(&Config{}, logger.NewLog()).Type()
	want := "slack"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestSlack_Notify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external test")
	}

	ctx := context.TODO()
	url := os.Getenv("WEBHOOK_URL")
	s := New(
		&Config{
			WebhookUrl:   url,
			ChannelNames: []string{"dev"},
			UserName:     "kompal-weather",
		},
		logger.NewLog(),
	)

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
