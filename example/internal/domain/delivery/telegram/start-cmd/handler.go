package start_cmd

import (
	"github.com/woyow/setupper/pkg/telegram"

	"github.com/sirupsen/logrus"
)

const (
	stateDefault = telegram.StateDefault
)

// Your business logic interface for this command
type service interface {
}

type Handler struct {
	service service
	tg      *telegram.Telegram
	log     *logrus.Logger
}

const (
	stateStartCommand = "/start"

	stateMenuCommand = "/menu"
)

func InitHandler(tg *telegram.Telegram, service service, log *logrus.Logger) {
	h := &Handler{
		service: service,
		tg:      tg,
		log:     log,
	}

	if err := h.tg.RegisterStates([]telegram.RegisterState{
		{
			stateStartCommand,
			h.handleStartCommand,
			true,
		},
	}); err != nil {
		panic(err)
	}
}

func (h *Handler) handleStartCommand(c *telegram.Context) telegram.StateFn {
	h.log.Debug("telegram delivery: handleStartCommand")

	return h.tg.SetStateAndCall(stateMenuCommand, c)
}
