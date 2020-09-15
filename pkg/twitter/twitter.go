package twitter

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/message"
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

// Notify message with images on Twitter
func (t Twitter) NotifyWithMedium(ctx context.Context, msg string, contents [][]byte) error {
	var medium []string
	for _, c := range contents {
		c := c
		m, err := t.uploadMedia(ctx, c)
		if err != nil {
			return fmt.Errorf("failed to upload media: %v", err)
		}
		medium = append(medium, m.MediaIDString)
	}
	// media_ids https://developer.twitter.com/en/docs/twitter-api/v1/tweets/post-and-engage/api-reference/post-statuses-update
	v := make(url.Values)
	v.Set("media_ids", strings.Join(medium, ","))

	if _, err := t.client.PostTweet(msg, v); err != nil {
		return fmt.Errorf("failed to tweet: %v", err)
	}
	return nil
}

func (t Twitter) uploadMedia(ctx context.Context, object []byte) (*anaconda.Media, error) {
	objStr := base64.StdEncoding.EncodeToString(object)
	m, err := t.client.UploadMedia(objStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get png object: %v", err)
	}

	return &m, nil
}
