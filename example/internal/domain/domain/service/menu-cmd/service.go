package menu_cmd

import (
	"context"

	commonEntity "github.com/woyow/setupper/example/internal/domain/domain/entity/common"
	entity "github.com/woyow/setupper/example/internal/domain/domain/entity/menu-cmd"
	types "github.com/woyow/setupper/example/internal/domain/domain/types/menu-cmd"

	"github.com/sirupsen/logrus"
	"github.com/woyow/setupper/pkg/marshaling/json"
	"github.com/woyow/setupper/pkg/telegram"
)

type api interface {
	SendMessage(dto commonEntity.SendMessageAPIDTO) error
	SendMessageWithInlineKeyboard(dto commonEntity.SendMessageWithInlineKeyboardAPIDTO) error
	EditMessageWithInlineKeyboard(dto commonEntity.EditMessageWithInlineKeyboardAPIDTO) error
	EditMessageText(dto commonEntity.EditMessageTextAPIDTO) error
	DeleteMessage(dto commonEntity.DeleteMessageAPIDTO) error
	DeleteMessages(dto commonEntity.DeleteMessagesAPIDTO) error
}


type repo interface {
	
}

type translate interface {
	Translate(key string, lang string) string
}

type Service struct {
	repo      repo
	api       api
	translate translate
	log       *logrus.Logger
}

func NewService(repo repo, api api, translate translate, log *logrus.Logger) *Service {
	s := &Service{
		repo:      repo,
		api:       api,
		translate: translate,
		log:       log,
	}

	return s
}

func (s *Service) MenuCommand(ctx context.Context, dto entity.MenuCommandServiceDTO) error {
	if err := s.sendMenu(sendMenuDTO{
		HandleCallback: telegram.HandleCallback{
			Lang:      dto.Lang,
			ChatID:    dto.ChatID,
			MessageID: 0,
		},
	}); err != nil {
		return err
	}

	return nil
}

const (
	callbackDataLoggingKey = "callback_data"
)

func (s *Service) MenuCommandCallback(ctx context.Context, dto entity.MenuCommandCallbackServiceDTO) (entity.MenuCommandCallbackOut, error) {
	var out entity.MenuCommandCallbackOut

	var callback entity.CallbackOut

	out.WaitCallback = true

	if err := json.Unmarshal([]byte(dto.CallbackData), &callback); err != nil {
		s.log.WithFields(logrus.Fields{callbackDataLoggingKey: dto.CallbackData}).Error("MenuCommandCallback - json.Unmarshal error: ", err)
		return out, nil
	}

	switch callback.Type {
	case types.ChooseMenu:
		out.WaitCallback = true

		if err := s.sendMenu(sendMenuDTO{
			HandleCallback: dto.HandleCallback,
		}); err != nil {
			return out, err
		}
	case types.ChooseMenuItem:
		out.WaitCallback = true

		var callback entity.ChooseMenuItemCallbackOut

		if err := json.Unmarshal([]byte(dto.CallbackData), &callback); err != nil {
			return out, err
		}

		if err := s.chooseMenuItem(ctx, chooseMenuItemDTO{
			callback:       callback,
			HandleCallback: dto.HandleCallback,
		}); err != nil {
			return out, err
		}
	
	default:
		s.log.WithFields(logrus.Fields{}).Error("MenuCommandCallback - unknown callback type: ", callback.Type)
	}

	return out, nil
}
