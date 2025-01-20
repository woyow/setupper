package telegram

import (
	"context"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/sirupsen/logrus"
)

type destructChatID int64

// chatInfo - need to implement for your telegram bot
type chatInfo interface {
	// IsBanned - is chat_id banned in your telegram bot
	IsBanned(ctx context.Context, dto IsBannedDTO) bool

	// CreateCurrentState - if the user entered the bot for the first time, you need to create default state in the database
	CreateCurrentState(ctx context.Context, dto CreateCurrentStateDTO) error

	// SetCurrentState - save chat state into database
	SetCurrentState(ctx context.Context, dto SetCurrentStateDTO) error

	// GetCurrentState - get current chat state from database
	GetCurrentState(ctx context.Context, dto GetCurrentStateDTO) (out GetCurrentStateOut, err error)
}

type bot struct {
	log        *logrus.Logger
	logFields  logrus.Fields
	chatInfo   chatInfo
	state      StateFn
	states     map[string]StateFn
	destructCh chan destructChatID
	chatID     int64
	tgBotName  string
}

type newBotDTO struct{
	chatInfo   chatInfo
	states     map[string]StateFn 
	destructCh chan destructChatID
	log        *logrus.Logger
	logFields  logrus.Fields
	tgBotName  string
}

func newBot(dto newBotDTO) func(chatID int64) echotron.Bot {
	return func(chatID int64) echotron.Bot {
		bot := &bot{
			log:        dto.log,
			logFields:  dto.logFields,
			chatInfo:   dto.chatInfo,
			state:      nil,
			states:     dto.states,
			destructCh: dto.destructCh,
			chatID:     chatID,
			tgBotName:  dto.tgBotName,
		}

		dto.logFields[chatIDLoggingKey] = chatID

		bot.setState()

		go bot.destruct()

		return bot
	}
}

func (b *bot) Update(update *echotron.Update) {
	c := Context(*update)
	b.state = b.state(&c)
}

func (b *bot) setState() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	isBanned := b.chatInfo.IsBanned(ctx, IsBannedDTO{
		ChatID: b.chatID,
	})
	if isBanned {
		b.state = b.states[StateBanned]
		return
	}

	chatCurrentState, err := b.chatInfo.GetCurrentState(ctx, GetCurrentStateDTO{
		ChatID: b.chatID,
	})
	if err != nil {
		switch err {
		case ErrChatCurrentStateNotExists:
			if err := b.chatInfo.CreateCurrentState(ctx, CreateCurrentStateDTO{
				State:  StateDefault,
				ChatID: b.chatID,
			}); err != nil {
				switch err {
				case ErrChatCurrentStateAlreadyExists:
					if b.log.Level == logrus.DebugLevel {
						b.log.WithFields(b.logFields).
							Error("setState - bot.stateService.CreateCurrentState error: ", err)
					}
				default:
					b.log.WithFields(b.logFields).
						Error("setState - bot.stateService.CreateCurrentState error: ", err)
				}
			}
		default:
			b.log.WithFields(b.logFields).
				Error("setState - bot.stateService.CreateCurrentState error: ", err)
		}

		b.state = b.states[StateDefault]
	} else {
		state, ok := b.states[chatCurrentState.State]
		if ok {
			b.state = state
			
			if b.log.Level == logrus.DebugLevel {
				b.log.WithFields(b.logFields).Debug("setState - Set " + chatCurrentState.State + " handler")
			}
		} else {
			b.state = b.states[StateDefault]

			if b.log.Level == logrus.DebugLevel {
				b.log.WithFields(b.logFields).Debug("setState - Set default handler")
			}
		}

		if chatCurrentState.State == StateBanned {
			b.state = b.states[StateDefault]
		}
	}
}

func (b *bot) destruct() {
	<-time.After(1 * time.Minute)

	b.destructCh <- destructChatID(b.chatID)
}
