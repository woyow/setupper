package common

import (
	"github.com/sirupsen/logrus"
	"github.com/woyow/setupper/pkg/telegram"
)

const (
	stateEmpty       = ""
	stateDefault     = telegram.StateDefault
	stateBanned      = telegram.StateBanned

	stateMenuCommand = "/menu"

	chatIDLoggingKey   = "chat_id"
	domainLoggingKey   = "domain"
	domainLoggingValue = "your-tg"
	layerLoggingKey    = "layer"
	layerLoggingValue  = "telegram delivery"
)

type service interface {
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
			stateEmpty,
			h.handleDefault,
			false,
		},
		{
			stateDefault,
			h.handleDefault,
			false,
		},
		{
			stateBanned,
			h.handleBanned,
			false,
		},
	}); err != nil {
		panic(err)
	}
}

func (h *Handler) handleDefault(c *telegram.Context) telegram.StateFn {
	if c.Message != nil {
		return h.handleMessage(c)
	}

	if c.CallbackQuery != nil {
		return h.handleCallbackQuery(c)
	}

	if c.EditedMessage != nil {
		return h.handleEditedMessage(c)
	}

	return h.tg.SetState(stateDefault, c)
}

// handleBanned - recursive set banned state
func (h *Handler) handleBanned(c *telegram.Context) telegram.StateFn {
	return h.tg.SetState(stateBanned, c)
}

func (h *Handler) handleMessage(c *telegram.Context) telegram.StateFn {
	if h.log.Level == logrus.DebugLevel {
		h.log.WithFields(logrus.Fields{
			chatIDLoggingKey: c.ChatID(),
			domainLoggingKey: domainLoggingValue,
			layerLoggingKey:  layerLoggingValue,
		}).Debug("handleMessage - Message text: ", c.Message.Text)
	
		h.log.WithFields(logrus.Fields{
			chatIDLoggingKey: c.ChatID(),
			domainLoggingKey: domainLoggingValue,
			layerLoggingKey:  layerLoggingValue,
		}).Debug("handleMessage - Message LanguageCode: ", c.Message.From.LanguageCode)
	}

	if h.tg.CheckCommand(c.Message.Text) {
		return h.tg.SetStateAndCall(c.Message.Text, c)
	}

	if c.Message.Contact != nil {
		return h.handleMessageContact(c)
	}

	h.log.Debug("handleMessage: unknown command")

	return h.tg.SetStateAndCall(stateMenuCommand, c)
}

func (h *Handler) handleMessageContact(c *telegram.Context) telegram.StateFn {
	h.log.Debug("handleMessageContact: ", 
		c.Message.Contact.FirstName, " ",
		c.Message.Contact.LastName, " ",
		c.Message.Contact.PhoneNumber,
	)

	return h.tg.SetStateAndCall(stateMenuCommand, c)
}

func (h *Handler) handleCallbackQuery(c *telegram.Context) telegram.StateFn {

	h.log.WithFields(logrus.Fields{
		chatIDLoggingKey: c.ChatID(),
		domainLoggingKey: domainLoggingValue,
		layerLoggingKey:  layerLoggingValue,
	}).Debug("handleCallbackQuery - CallbackQuery Data: ", c.CallbackQuery.Data)

	h.log.WithFields(logrus.Fields{
		chatIDLoggingKey: c.ChatID(),
		domainLoggingKey: domainLoggingValue,
		layerLoggingKey:  layerLoggingValue,
	}).Debug("handleCallbackQuery - CallbackQuery Message.ID: ", c.CallbackQuery.Message.ID)

	return h.tg.SetState(stateDefault, c)
}

func (h *Handler) handleEditedMessage(c *telegram.Context) telegram.StateFn {

	h.log.WithFields(logrus.Fields{
		chatIDLoggingKey: c.ChatID(),
		domainLoggingKey: domainLoggingValue,
		layerLoggingKey:  layerLoggingValue,
	}).Debug("handleEditedMessage - EditedMessage text: ", c.EditedMessage.Text)

	h.log.WithFields(logrus.Fields{
		chatIDLoggingKey: c.ChatID(),
		domainLoggingKey: domainLoggingValue,
		layerLoggingKey:  layerLoggingValue,
	}).Debug("handleEditedMessage - EditedMessage ID: ", c.EditedMessage.ID)

	return h.tg.SetState(stateDefault, c)
}
