package echotron

import (
	"os"

	"github.com/NicoNex/echotron/v3"
	"github.com/sirupsen/logrus"
)

const (
	setupLoggingKey   = "setup"
	setupLoggingValue = "echotron"
)

type Echotron struct {
	api         *echotron.API
	token       string
	webhookURL  string
	webhookAddr string
}

func NewEchotron(cfg *Config, log *logrus.Logger) (*Echotron, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	token := os.Getenv(cfg.TokenEnvKey)
	webhookURL := os.Getenv(cfg.WebhookURLEnvKey)
	webhookAddr := os.Getenv(cfg.WebhookHTTPAddrEnvKey)

	if webhookURL != "" {
		if webhookAddr == "" {
			return nil, ErrEmptyWebhookAddress
		}
	}

	api := echotron.NewAPI(token)

	if log.Level == logrus.DebugLevel {
		log.WithField(setupLoggingKey, setupLoggingValue).Debug("NewEchotron - API has been initialized")
	}

	return &Echotron{
		api:         &api,
		token:       token,
		webhookURL:  webhookURL,
		webhookAddr: webhookAddr,
	}, nil
}

func (e *Echotron) GetToken() string {
	return e.token
}

func (e *Echotron) GetAPI() *echotron.API {
	return e.api
}

func (e *Echotron) GetWebhookURL() string {
	return e.webhookURL
}

func (e *Echotron) GetWebhookAddr() string {
	return e.webhookAddr
}
