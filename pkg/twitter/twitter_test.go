package twitter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/status"
)

func TestTwitter_Type(t *testing.T) {
	got := New(&Config{}, logger.NewLog()).Type()
	want := "twitter"
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestTwitter_Notify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external test")
	}

	ctx := context.TODO()
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiKeySecret := os.Getenv("TWITTER_API_KEY_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	tw := New(
		&Config{
			APIKey:            apiKey,
			APIKeySecret:      apiKeySecret,
			AccessToken:       accessToken,
			AccessTokenSecret: accessTokenSecret,
		},
		logger.NewLog(),
	)

	err := tw.Notify(ctx, &analyzer.Result{
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
