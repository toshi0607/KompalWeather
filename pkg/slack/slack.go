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

// Slack is representation of Slack
type Slack struct {
	config *Config
	log    logger.Logger
}

// Config is a configuration of Slack
type Config struct {
	WebhookURL   string
	ChannelNames []string
	UserName     string
}

// New builds new Slack
func New(config *Config, log logger.Logger) *Slack {
	return &Slack{
		config: config,
		log:    log,
	}
}

// Type returns the type of the notifier
func (s Slack) Type() string {
	return "slack"
}

// Notify notifies result in Slack channel
func (s Slack) Notify(ctx context.Context, result *analyzer.Result) error {
	if result.MaleTrend == analyzer.Constant && result.FemaleTrend == analyzer.Constant {
		s.log.Info("skip slack notification")
		return nil
	}

	m := message.Build(result)
	j := `{"channel":"` + s.config.ChannelNames[0] + `","username":"` + s.config.UserName + `","text":"` + m + `"}`
	s.log.Info(j)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.config.WebhookURL,
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.log.Error("failed to close resp body: %s", err)
		}
	}()

	return nil
}

func (s Slack) NotifyWithMedium(ctx context.Context, msg string, contents [][]byte) error {
	s.log.Info("skip slack notification")
	return nil
}
