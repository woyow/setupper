package echotron

type Config struct {
	BotNameEnvKey         string `yaml:"bot_name_env_key" default:"TG_BOT_NAME"`
	TokenEnvKey           string `yaml:"token_env_key" default:"TG_TOKEN"`
	WebhookURLEnvKey      string `yaml:"webhook_url_env_key" default:"TG_WEBHOOK_URL"`
	WebhookHTTPAddrEnvKey string `yaml:"webhook_http_addr_env_key" default:"TG_WEBHOOK_HTTP_ADDRES"`
}

func (c *Config) Validate() error {
	if c.BotNameEnvKey == "" {
		return ErrEmptyBotNameEnvKey
	}
	if c.TokenEnvKey == "" {
		return ErrEmptyTokenEnvKey
	}
	if c.WebhookURLEnvKey == "" {
		return ErrEmptyWebhookURLEnvKey
	}
	if c.WebhookHTTPAddrEnvKey == "" {
		return ErrEmptyWebhookHTTPAddrEnvKey
	}
	return nil
}