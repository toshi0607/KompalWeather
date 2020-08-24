package slack

import (
	"bytes"
	"context"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/message"
)

type Slack struct {
	config *Config
}

type Config struct {
	WebhookUrl   string
	ChannelNames []string
	UserName     string
	//Frequency    config.Frequency
}

func New(config *Config) *Slack {

	return &Slack{
		config: config,
	}
}

func (s Slack) Notify(ctx context.Context, result *analyzer.Result) error {
	m := message.Build(result)
	j := `{"channel":"` + s.config.ChannelNames[0] + `","username":"` + s.config.UserName + `","text":"` + m + `"}`
	req, err := http.NewRequest(
		http.MethodPost,
		s.config.WebhookUrl,
		bytes.NewBuffer([]byte(j)),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
