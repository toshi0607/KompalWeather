package config

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/slack"
	"github.com/toshi0607/kompal-weather/pkg/storage"
)

type Frequency int

const (
	Unknown    = 0
	Increasing = 1
	Decreasing = 2
	All        = 3
)

type Config struct {
	Slack        *slack.Config
	Kompal       *kompal.Config
	Sheets       *storage.SheetsConfig
	ServerPort   int
	GCPProjectID string

	secret *secret.Secret
}

type env struct {
	GCPProjectID              string   `envconfig:"GCP_PROJECT_ID" required:"true"`
	ServerPort                int      `envconfig:"SERVER_PORT" required:"true"`
	SlackChannelNames         []string `envconfig:"SLACK_CHANNEL_NAMES" required:"true"`
	SlackUserName             string   `envconfig:"SLACK_USER_NAME" required:"true" default:"kompal-weather"`
	SlackWebhookUrlSecretName string   `envconfig:"SLACK_WEBHOOK_URL_SECRET_NAME" required:"true"`
	KompalUrlSecretName       string   `envconfig:"KOMPAL_URL_SECRET_NAME" required:"true"`
	SpreadSheetID             string   `envconfig:"SPREAD_SHEET_ID" required:"true"`
	SheetID                   uint     `envconfig:"SHEET_ID" required:"true"`
}

func New(secret *secret.Secret) *Config {
	return &Config{
		secret: secret,
	}
}

func (c *Config) Init() error {
	ctx := context.TODO()
	var e env
	if err := envconfig.Process("", &e); err != nil {
		return fmt.Errorf("failed to process envconfig: %s", err)
	}

	// Secret
	kompalUrl, err := c.secret.Get(ctx, e.KompalUrlSecretName)
	if err != nil {
		return err
	}
	slackWebhookUrl, err := c.secret.Get(ctx, e.SlackWebhookUrlSecretName)
	if err != nil {
		return err
	}

	c.Kompal.URL = kompalUrl
	c.Slack = &slack.Config{
		WebhookUrl:   slackWebhookUrl,
		UserName:     e.SlackUserName,
		ChannelNames: e.SlackChannelNames,
	}
	c.GCPProjectID = e.GCPProjectID
	c.ServerPort = e.ServerPort
	c.Sheets = &storage.SheetsConfig{
		SpreadSheetID: e.SpreadSheetID,
		SheetId:       e.SheetID,
	}

	return nil
}
