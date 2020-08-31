package twitter

import (
	"context"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/message"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

// Twitter is representation of Twitter
type Twitter struct {
	client *anaconda.TwitterApi
	log    logger.Logger
}

// Config is a configuration of Twitter
type Config struct {
	APIKey            string
	APIKeySecret      string
	AccessToken       string
	AccessTokenSecret string
}

// New builds new Twitter
func New(config *Config, log logger.Logger) *Twitter {
	return &Twitter{
		client: anaconda.NewTwitterApiWithCredentials(
			config.AccessToken,
			config.AccessTokenSecret,
			config.APIKey,
			config.APIKeySecret,
		),
		log: log,
	}
}

// Type returns the type of the notifier
func (t Twitter) Type() string {
	return "twitter"
}

// Notify notifies result on Twitter
func (t Twitter) Notify(ctx context.Context, result *analyzer.Result) error {
	if result.MaleTrend == analyzer.Constant && result.FemaleTrend == analyzer.Constant {
		t.log.Info("skip twitter notification")
		return nil
	}

	m := message.Build(result)
	if _, err := t.client.PostTweet(m, nil); err != nil {
		return fmt.Errorf("failed to tweet: %v", err)
	}
	return nil
}
