package twitter

import (
	"context"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/toshi0607/kompal-weather/pkg/message"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
)

type Twitter struct {
	client *anaconda.TwitterApi
}

type Config struct {
	APIKey            string
	APIKeySecret      string
	AccessToken       string
	AccessTokenSecret string
}

func New(config *Config) *Twitter {
	return &Twitter{
		client: anaconda.NewTwitterApiWithCredentials(
			config.AccessToken,
			config.AccessTokenSecret,
			config.APIKey,
			config.APIKeySecret,
		),
	}
}

func (t Twitter) Type() string {
	return "twitter"
}

func (t Twitter) Notify(ctx context.Context, result *analyzer.Result) error {
	if result.MaleTrend == analyzer.Constant && result.FemaleTrend == analyzer.Constant {
		fmt.Print("skip twitter notification\n")
		return nil
	}

	m := message.Build(result)
	if _, err := t.client.PostTweet(m, nil); err != nil {
		return fmt.Errorf("failed to tweet: %v", err)
	}
	return nil
}