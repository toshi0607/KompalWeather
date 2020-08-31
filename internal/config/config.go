package config

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/secret"
	"github.com/toshi0607/kompal-weather/pkg/slack"
	"github.com/toshi0607/kompal-weather/pkg/storage"
	"github.com/toshi0607/kompal-weather/pkg/twitter"
)

// Config is a whole configuration for komal-weather application
type Config struct {
	Slack        *slack.Config
	Twitter      *twitter.Config
	Kompal       *kompal.Config
	Sheets       *storage.SheetsConfig
	ServerPort   int
	GCPProjectID string
	Version      string
	ServiceName  string
	Environment  string

	secret *secret.Secret
}

type env struct {
	// GCPProjectID is a GCP project id where this application is hosted
	GCPProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
	// ServerPort is a port number this application listens to
	ServerPort int `envconfig:"SERVER_PORT" required:"true"`
	// SlackChannelNames are names of notification target Slack channels
	SlackChannelNames []string `envconfig:"SLACK_CHANNEL_NAMES" required:"true"`
	// SlackUserName is a user name used when notifying slack channels
	SlackUserName string `envconfig:"SLACK_USER_NAME" required:"true" default:"kompal-weather"`
	// SlackWebhookURLSecretName is a secret_id of Slack incoming webhook URL
	SlackWebhookURLSecretName string `envconfig:"SLACK_WEBHOOK_URL_SECRET_NAME" required:"true"`
	// TwitterAPIKeySecretName is a secret_id of Twitter API key
	TwitterAPIKeySecretName string `envconfig:"TWITTER_API_KEY_SECRET_NAME" required:"false"`
	// TwitterAPIKeySecretSecretName is a secret_id of Twitter API key secret
	TwitterAPIKeySecretSecretName string `envconfig:"TWITTER_API_KEY_SECRET_SECRET_NAME" required:"false"`
	// TwitterAccessTokenSecretName is a secret_id of Twitter access token
	TwitterAccessTokenSecretName string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET_NAME" required:"false"`
	// TwitterAccessTokenSecretSecretName is a secret_id of Twitter access token secret
	TwitterAccessTokenSecretSecretName string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET_SECRET_NAME" required:"false"`
	// KompalURLSecretName is a secret_id of API endpoint for Komparu-yu
	KompalURLSecretName string `envconfig:"KOMPAL_URL_SECRET_NAME" required:"true"`
	// SpreadSheetID is a id of spreadheet
	SpreadSheetID string `envconfig:"SPREAD_SHEET_ID" required:"true"`
	// SheetID is a id of each sheet
	SheetID uint `envconfig:"SHEET_ID" required:"true"`
	// ServiceName is a name of this service
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`
	// Version is a version of this application
	Version string `envconfig:"VERSION" required:"true"`
	// Environment is environment of current application. env const must be used.
	Environment string `envconfig:"ENVIRONMENT" required:"true"`
}

const (
	envLocal       = "local"
	envDevelopment = "development"
	envProduction  = "production"
)

func (e *env) validate() error {
	if e.Environment != envDevelopment && e.Environment != envProduction && e.Environment != envLocal {
		return fmt.Errorf("env is invalid, env:%s", e.Environment)
	}
	return nil
}

// New builds new Config
func New(secret *secret.Secret) *Config {
	return &Config{
		Slack:  &slack.Config{},
		Kompal: &kompal.Config{},
		Sheets: &storage.SheetsConfig{},
		secret: secret,
	}
}

// Init inits Config
func (c *Config) Init() error {
	ctx := context.TODO()

	// Environment variable
	var e env
	if err := envconfig.Process("", &e); err != nil {
		return fmt.Errorf("failed to process envconfig: %s", err)
	}
	if err := e.validate(); err != nil {
		return fmt.Errorf("failed to validate env: %s", err)
	}

	// Secret
	c.secret.AddGCPProjectID(e.GCPProjectID)
	kompalURL, err := c.secret.Get(ctx, e.KompalURLSecretName)
	if err != nil {
		return err
	}
	slackWebhookURL, err := c.secret.Get(ctx, e.SlackWebhookURLSecretName)
	if err != nil {
		return err
	}
	accessToken, err := c.secret.Get(ctx, e.TwitterAccessTokenSecretName)
	if err != nil {
		return err
	}
	accessTokenSecret, err := c.secret.Get(ctx, e.TwitterAccessTokenSecretSecretName)
	if err != nil {
		return err
	}
	apiKey, err := c.secret.Get(ctx, e.TwitterAPIKeySecretName)
	if err != nil {
		return err
	}
	apiKeySecret, err := c.secret.Get(ctx, e.TwitterAPIKeySecretSecretName)
	if err != nil {
		return err
	}

	c.Kompal.URL = kompalURL
	c.Slack = &slack.Config{
		WebhookURL:   slackWebhookURL,
		UserName:     e.SlackUserName,
		ChannelNames: e.SlackChannelNames,
	}
	c.Twitter = &twitter.Config{
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
		APIKey:            apiKey,
		APIKeySecret:      apiKeySecret,
	}
	c.GCPProjectID = e.GCPProjectID
	c.ServerPort = e.ServerPort
	c.ServiceName = e.ServiceName
	c.Version = e.Version
	c.Environment = e.Environment
	c.Sheets = &storage.SheetsConfig{
		SpreadSheetID: e.SpreadSheetID,
		SheetID:       e.SheetID,
	}

	return nil
}

// IsLocal returns env is local or not
func (c *Config) IsLocal() bool {
	return c.Environment == envLocal
}
