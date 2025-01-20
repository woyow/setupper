package echotron

import "errors"

var (
	ErrEmptyBotNameEnvKey         = errors.New("empty bot_name_env_key")
	ErrEmptyTokenEnvKey           = errors.New("empty token_env_key")
	ErrEmptyWebhookURLEnvKey      = errors.New("empty webhook_url_env_key")
	ErrEmptyWebhookHTTPAddrEnvKey = errors.New("empty webhook_http_addr_env_key")

	ErrEmptyWebhookAddress = errors.New("empty webhook address")
)