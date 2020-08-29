package slack

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/message"
)

type Slack struct {
	config *Config
	log    logger.Logger
}

type Config struct {
	WebhookUrl   string
	ChannelNames []string
	UserName     string
	//Frequency    config.Frequency
}

func New(config *Config, log logger.Logger) *Slack {
	return &Slack{
		config: config,
		log:    log,
	}
}

func (s Slack) Type() string {
	return "slack"
}

func (s Slack) Notify(ctx context.Context, result *analyzer.Result) error {
	if result.MaleTrend == analyzer.Constant && result.FemaleTrend == analyzer.Constant {
		s.log.Info("skip slack notification")
		return nil
	}

	m := message.Build(result)
	j := `{"channel":"` + s.config.ChannelNames[0] + `","username":"` + s.config.UserName + `","text":"` + m + `"}`
	s.log.Info(j)
	req, err := http.NewRequest(
		http.MethodPost,
		s.config.WebhookUrl,
		bytes.NewBuffer([]byte(j)),
	)
	if err != nil {
		return fmt.Errorf("failed to create slack request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to slack: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
