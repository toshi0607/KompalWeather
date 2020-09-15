package config

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/toshi0607/kompal-weather/pkg/secret"
)

type VisualizerConfig struct {
	Mail string
	PW   string
	// Common config
	BucketName   string
	ServerPort   int
	GCPProjectID string
	Version      string
	ServiceName  string
	Environment  string

	secret *secret.Secret
}

type visualizerEnv struct {
	BucketName     string `envconfig:"GCS_BUCKET_NAME" required:"true"`
	MailSecretName string `envconfig:"MAIL_SECRET_NAME" required:"true"`
	PWSecretName   string `envconfig:"PW_SECRET_NAME" required:"true"`
	// GCPProjectID is a GCP project id where this application is hosted
	GCPProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
	// ServerPort is a port number this application listens to
	ServerPort int `envconfig:"SERVER_PORT" required:"true"`
	// SlackChannelNames are names of notification target Slack channels
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`
	// Version is a version of this application
	Version string `envconfig:"VERSION" required:"true"`
	// Environment is environment of current application. env const must be used.
	Environment string `envconfig:"ENVIRONMENT" required:"true"`
}

func (e *visualizerEnv) validate() error {
	if e.Environment != envDevelopment && e.Environment != envProduction && e.Environment != envLocal {
		return fmt.Errorf("coreEnv is invalid, coreEnv:%s", e.Environment)
	}
	return nil
}

// NewCore builds new CoreConfig
func NewVisualizer(secret *secret.Secret) *VisualizerConfig {
	return &VisualizerConfig{
		secret: secret,
	}
}

func (c *VisualizerConfig) Init() error {
	ctx := context.TODO()

	// Environment variable
	var e visualizerEnv
	if err := envconfig.Process("", &e); err != nil {
		return fmt.Errorf("failed to process envconfig: %s", err)
	}
	if err := e.validate(); err != nil {
		return fmt.Errorf("failed to validate coreEnv: %s", err)
	}

	// Secret
	c.secret.AddGCPProjectID(e.GCPProjectID)
	mail, err := c.secret.Get(ctx, e.MailSecretName)
	if err != nil {
		return err
	}
	pw, err := c.secret.Get(ctx, e.PWSecretName)
	if err != nil {
		return err
	}

	c.Mail = mail
	c.PW = pw
	c.BucketName = e.BucketName
	c.GCPProjectID = e.GCPProjectID
	c.ServerPort = e.ServerPort
	c.ServiceName = e.ServiceName
	c.Version = e.Version
	c.Environment = e.Environment
	return nil
}

// IsLocal returns coreEnv is local or not
func (c *VisualizerConfig) IsLocal() bool {
	return c.Environment == envLocal
}
