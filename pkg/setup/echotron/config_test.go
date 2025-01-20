package echotron

import (
	"testing"
)


func TestValidate(t *testing.T) {
	testCases := []struct{
		name  string
		cfg   Config
		expErr error
	}{
		{
			name: "empty env key",
			cfg: Config{
				BotNameEnvKey: "",
				TokenEnvKey: "TG_TOKEN",
				WebhookURLEnvKey: "TG_WEBHOOK_URL",
				WebhookHTTPAddrEnvKey: "TG_WEBHOOK_ADDRESS",
			},
			expErr: ErrEmptyBotNameEnvKey,
		},
		{
			name: "empty token name",
			cfg: Config{
				BotNameEnvKey: "your_tg",
				TokenEnvKey: "",
				WebhookURLEnvKey: "TG_WEBHOOK_URL",
				WebhookHTTPAddrEnvKey: "TG_WEBHOOK_ADDRESS",
			},
			expErr: ErrEmptyTokenEnvKey,
		},
		{
			name: "empty webhook url",
			cfg: Config{
				BotNameEnvKey: "your_tg",
				TokenEnvKey: "TG_TOKEN",
				WebhookURLEnvKey: "",
				WebhookHTTPAddrEnvKey: "TG_WEBHOOK_ADDRESS",
			},
			expErr: ErrEmptyWebhookURLEnvKey,
		},
		{
			name: "empty webhook url",
			cfg: Config{
				BotNameEnvKey: "your_tg",
				TokenEnvKey: "TG_TOKEN",
				WebhookURLEnvKey: "TG_WEBHOOK_URL",
				WebhookHTTPAddrEnvKey: "",
			},
			expErr: ErrEmptyWebhookHTTPAddrEnvKey,
		},
	}

	for _, testCase := range testCases {
		if err := testCase.cfg.Validate(); err != testCase.expErr {
			t.Errorf("got: %v, want: %v", err, testCase.expErr)
		}
	}
}