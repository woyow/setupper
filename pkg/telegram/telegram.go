package telegram

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	setupEchotron "github.com/woyow/setupper/pkg/setup/echotron"

	"github.com/NicoNex/echotron/v3"
	"github.com/sirupsen/logrus"
)

const (
	chatIDLoggingKey     = "chat_id"
	tgBotNameLoggingKey  = "tg_bot_name"
	layerLoggingKey      = "layer"
	layerLoggingValue    = "delivery"
)

type StateFn func(ctx *Context) StateFn

type Telegram struct {
	chatInfo     chatInfo
	mu           *sync.RWMutex
	destructCh   chan destructChatID
	states       map[string]StateFn
	commands     map[string]struct{}
	log          *logrus.Logger
	logFields    logrus.Fields
	tgBotName    string
	token        string
	webhookURL   string
	webhookAddr  string
}

func NewTelegram(setupEchotron *setupEchotron.Echotron, chatInfo chatInfo, log *logrus.Logger, tgBotName string) *Telegram {
	return &Telegram{
		chatInfo:     chatInfo,
		mu:           &sync.RWMutex{},
		destructCh:   make(chan destructChatID),
		states:       make(map[string]StateFn),
		commands:     make(map[string]struct{}),
		log:          log,
		logFields: logrus.Fields{
			tgBotNameLoggingKey: tgBotName,
			layerLoggingKey:  layerLoggingValue,
		},
		tgBotName:      tgBotName,
		token:       setupEchotron.GetToken(),
		webhookURL:  setupEchotron.GetWebhookURL(),
		webhookAddr: setupEchotron.GetWebhookAddr(),
	}
}

func (t *Telegram) CheckCommand(command string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.commands[command]
	return ok
}

func (t *Telegram) destructBot(dispatcher *echotron.Dispatcher) {
	for {
		select {
		case b := <-t.destructCh:
			dispatcher.DelSession(int64(b))
			
			if t.log.Level == logrus.DebugLevel {
				t.log.WithFields(logrus.Fields{
					tgBotNameLoggingKey: t.tgBotName,
					layerLoggingKey:  layerLoggingValue,
				}).Info("telegram: destructBot - destruct bot with chatID: ", int64(b))
			}
		}
	}
}

type RegisterState struct {
	Name      string
	StateFn   StateFn
	IsCommand bool
}

func (t *Telegram) RegisterStates(states []RegisterState) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := range states {
		if _, ok := t.states[states[i].Name]; ok {
			return fmt.Errorf("telegram handler %s is already registered", states[i].Name)
		}

		t.states[states[i].Name] = states[i].StateFn

		if states[i].IsCommand {
			t.commands[states[i].Name] = struct{}{}
		}
	}

	return nil
}

func (t *Telegram) SetStateAndCall(state string, c *Context) StateFn {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stateFn, ok := t.states[state]
	if !ok {
		t.log.WithFields(logrus.Fields{
			tgBotNameLoggingKey: t.tgBotName,
			chatIDLoggingKey: c.ChatID(),
		}).Error("telegram: SetStateAndCall - Unknown state ", state)

		return t.SetStateAndCall(StateDefault, c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := t.chatInfo.SetCurrentState(ctx, SetCurrentStateDTO{
		State:  state,
		ChatID: c.ChatID(),
	}); err != nil {
		t.log.WithFields(logrus.Fields{
			tgBotNameLoggingKey: t.tgBotName,
			chatIDLoggingKey: c.ChatID(),
		}).Error("telegram: SetStateAndCall - t.chatInfo.SetCurrentState error: ", err)
	}

	return stateFn(c)
}

func (t *Telegram) SetState(state string, c *Context) StateFn {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stateFn, ok := t.states[state]
	if !ok {
		t.log.WithFields(logrus.Fields{
			tgBotNameLoggingKey: t.tgBotName,
			chatIDLoggingKey: c.ChatID(),
		}).Error("telegram: SetStateAndCall - Unknown state ", state)

		return t.SetState(StateDefault, c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := t.chatInfo.SetCurrentState(ctx, SetCurrentStateDTO{
		State:  state,
		ChatID: c.ChatID(),
	}); err != nil {
		t.log.WithFields(logrus.Fields{
			chatIDLoggingKey: c.ChatID(),
			tgBotNameLoggingKey: t.tgBotName,
		}).Error("telegram: SetState - t.chatInfo.SetCurrentState error: ", err)
	}

	return stateFn
}

func (t *Telegram) Run() error {
	if t.log.Level == logrus.DebugLevel {
		defer t.log.WithFields(t.logFields).Debug("telegram: Run - stop telegram bot")
	}

	dispatcher := echotron.NewDispatcher(t.token, newBot(newBotDTO{
		states:     t.states, 
		destructCh: t.destructCh, 
		chatInfo:   t.chatInfo, 
		log:        t.log, 
		logFields:  t.logFields, 
		tgBotName:  t.tgBotName,
	}))

	go t.destructBot(dispatcher)

	if t.log.Level == logrus.DebugLevel {
		t.log.WithFields(t.logFields).Debug("telegram: Run - start telegram bot")
	}

	if t.webhookURL != "" {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("OK\n"))
			if err != nil {
				t.log.WithFields(logrus.Fields{
					tgBotNameLoggingKey: t.tgBotName,
				}).Error("telegram: Run - w.Write error", err)
			}
			w.WriteHeader(http.StatusOK)
		})

		dispatcher.SetHTTPServer(
			&http.Server{
				Addr:                         t.webhookAddr,
				Handler:                      mux,
				DisableGeneralOptionsHandler: false,
				TLSConfig:                    nil,
				ReadTimeout:                  0,
				ReadHeaderTimeout:            0,
				WriteTimeout:                 0,
				IdleTimeout:                  0,
				MaxHeaderBytes:               0,
				TLSNextProto:                 nil,
				ConnState:                    nil,
				ErrorLog:                     nil,
				BaseContext:                  nil,
				ConnContext:                  nil,
			})

		return dispatcher.ListenWebhookOptions(t.webhookURL, true, &echotron.WebhookOptions{
			IPAddress:   "",
			SecretToken: "",
			Certificate: echotron.InputFile{},
			AllowedUpdates: []echotron.UpdateType{
				echotron.MessageUpdate,
				echotron.EditedMessageUpdate,
				echotron.ChannelPostUpdate,
				echotron.EditedChannelPostUpdate,
				echotron.InlineQueryUpdate,
				echotron.ChosenInlineResultUpdate,
				echotron.CallbackQueryUpdate,
				echotron.ShippingQueryUpdate,
				echotron.PreCheckoutQueryUpdate,
				echotron.MyChatMemberUpdate,
				echotron.ChatMemberUpdate,
			},
			MaxConnections: 10,
		})
	} else {
		return dispatcher.Poll()
	}
}
