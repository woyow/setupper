package menu_cmd

import (
	"context"
	"time"

	entity "github.com/woyow/setupper/example/internal/domain/domain/entity/menu-cmd"

	"github.com/sirupsen/logrus"
	"github.com/woyow/setupper/pkg/telegram"
)

const (
	stateDefault = telegram.StateDefault
)

type service interface {
	MenuCommand(ctx context.Context, dto entity.MenuCommandServiceDTO) error
	MenuCommandCallback(ctx context.Context, dto entity.MenuCommandCallbackServiceDTO) (entity.MenuCommandCallbackOut, error)
}

type Handler struct {
	service service
	tg      *telegram.Telegram
	log     *logrus.Logger
}

const (
	stateMenuCommand         = "/menu"
	stateMenuCommandCallback = "/menu_callback"
)

func InitHandler(tg *telegram.Telegram, service service, log *logrus.Logger) {
	h := &Handler{
		service: service,
		tg:      tg,
		log:     log,
	}

	if err := h.tg.RegisterStates([]telegram.RegisterState{
		{
			stateMenuCommand,
			h.handleMenuCommand,
			true,
		},
		{
			stateMenuCommandCallback,
			h.handleMenuCommandCallback,
			false,
		},
	}); err != nil {
		panic(err)
	}
}

func (h *Handler) handleMenuCommand(c *telegram.Context) telegram.StateFn {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.service.MenuCommand(ctx, entity.MenuCommandServiceDTO{
		HandleCommand: telegram.HandleCommand{
			Lang:   c.Message.From.LanguageCode,
			ChatID: c.ChatID(),
		}},
	); err != nil {
		return h.tg.SetState(stateDefault, c)
	}

	return h.tg.SetState(stateMenuCommandCallback, c)
}

func (h *Handler) handleMenuCommandCallback(c *telegram.Context) telegram.StateFn {
	if c.Message != nil {
		if h.tg.CheckCommand(c.Message.Text) {
			return h.tg.SetStateAndCall(c.Message.Text, c)
		} else {
			return h.tg.SetState(stateMenuCommandCallback, c)
		}
	}

	if c.CallbackQuery == nil {
		return h.tg.SetState(stateMenuCommandCallback, c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	out, err := h.service.MenuCommandCallback(ctx, entity.MenuCommandCallbackServiceDTO{
		HandleCallback: telegram.HandleCallback{
			Lang:      c.CallbackQuery.From.LanguageCode,
			ChatID:    c.ChatID(),
			MessageID: c.CallbackQuery.Message.ID,
		},
		CallbackData: c.CallbackQuery.Data,
	})
	if err != nil {
		return h.tg.SetState(stateMenuCommandCallback, c)
	}

	if out.WaitCallback {
		return h.tg.SetState(stateMenuCommandCallback, c)		
	} else {
		return h.tg.SetState(stateMenuCommandCallback, c)
	}
}
