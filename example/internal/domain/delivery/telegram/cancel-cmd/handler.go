package cancel_cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/woyow/setupper/pkg/telegram"
)

const (
	stateDefault = telegram.StateDefault

	stateCancelCommand   = "/cancel"
	stateMenuCommand     = "/menu"
)

// Your business logic interface for this command
type service interface{
}

type Handler struct {
	service service
	tg      *telegram.Telegram
	log     *logrus.Logger
}


func InitHandler(tg *telegram.Telegram, service service, log *logrus.Logger) {
	h := &Handler{
		service: service,
		tg:      tg,
		log:     log,
	}

	if err := h.tg.RegisterStates([]telegram.RegisterState{
		{
			stateCancelCommand,
			h.handleCancelCommand,
			true,
		},
	}); err != nil {
		panic(err)
	}
}

func (h *Handler) handleCancelCommand(c *telegram.Context) telegram.StateFn {
	h.log.Debug("telegram delivery: handleCancelCommand")

	return h.tg.SetStateAndCall(stateMenuCommand, c)
}
