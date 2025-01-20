package tgapi

import (
	"github.com/NicoNex/echotron/v3"
	"github.com/sirupsen/logrus"
)

const (
	textLoggingKey     = "text"
	chatIDLoggingKey   = "chat_id"
	tgBotNameLoggingKey   = "tg_bot_name"
	tgBotNameLoggingValue = "your-tg"
	infraLoggingKey    = "infra"
	infraLoggingValue  = "tgapi"
	fileIDLoggingKey   = "file_id"
)

type API struct {
	api *echotron.API
	log *logrus.Logger
}

func NewAPI(api *echotron.API, log *logrus.Logger) *API {
	return &API{
		api: api,
		log: log,
	}
}
