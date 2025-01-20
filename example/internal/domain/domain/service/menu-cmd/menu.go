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

const (
	menuTranslateKey          = "menu"
	menuItemTest1TranslateKey = "menu_item_test1"
	menuItemTest2TranslateKey = "menu_item_test2"
	menuItemTest3TranslateKey = "menu_item_test3"
)

type sendMenuDTO struct {
	telegram.HandleCallback
}

func (s *Service) sendMenu(dto sendMenuDTO) error {
	buttons := s.sendMenuButtons(sendMenuButtonsDTO{
		lang: dto.Lang,
	})

	text := s.translate.Translate(menuTranslateKey, dto.Lang)

	if dto.MessageID == 0 {
		if err := s.api.SendMessageWithInlineKeyboard(commonEntity.SendMessageWithInlineKeyboardAPIDTO{
			Buttons: buttons,
			Text:    text,
			ChatID:  dto.ChatID,
			Options: nil,
		}); err != nil {
			s.log.WithFields(logrus.Fields{}).Error("sendMenu - s.api.SendMessageWithInlineKeyboard error: ", err)
			return err
		}
	} else {
		if err := s.api.EditMessageText(commonEntity.EditMessageTextAPIDTO{
			Buttons:   buttons,
			Options:   nil,
			Text:      text,
			ChatID:    dto.ChatID,
			MessageID: dto.MessageID,
		}); err != nil {
			s.log.WithFields(logrus.Fields{}).Error("sendMenu - s.api.EditMessageText error: ", err)
			return err
		}
	}
	return nil
}

type sendMenuButtonsDTO struct {
	lang string
}

func (s *Service) sendMenuButtons(dto sendMenuButtonsDTO) [][]commonEntity.InlineButton {
	eb := [][]types.MenuItem{
		{
			types.MenuItemTest1,
		},
		{
			types.MenuItemTest2,
			types.MenuItemTest3,
		},
	}

	buttons := make([][]commonEntity.InlineButton, 0, len(eb))

	for _, item := range eb {
		buttons = append(buttons, func() []commonEntity.InlineButton {
			buttonsLine := make([]commonEntity.InlineButton, 0, len(item))
			for _, item2 := range item {
				buttonsLine = append(buttonsLine, commonEntity.InlineButton{
					Text: func() string {
						var translateKey string
						var postFix string

						switch item2 {
						case types.MenuItemTest1:
							translateKey = menuItemTest1TranslateKey
						case types.MenuItemTest2:
							translateKey = menuItemTest2TranslateKey
						case types.MenuItemTest3:
							translateKey = menuItemTest3TranslateKey
						}
						return s.translate.Translate(translateKey, dto.lang) + postFix
					}(),
					CallbackData: func() string {
						callback := entity.ChooseMenuItemCallbackOut{
							CallbackOut: entity.CallbackOut{
								Type: types.ChooseMenuItem,
							},
							MenuItem: item2,
						}

						callbackData, err := json.Marshal(&callback)
						if err != nil {
							s.log.WithFields(logrus.Fields{}).Error("sendMenuButtons - json.Marshal error: ", err.Error())
						}

						return string(callbackData)
					}(),
					URL: "",
				})
			}
			return buttonsLine
		}())
	}


	return buttons
}

type chooseMenuItemDTO struct {
	callback entity.ChooseMenuItemCallbackOut
	telegram.HandleCallback
}

func (s *Service) chooseMenuItem(_ context.Context, dto chooseMenuItemDTO) error {
	switch dto.callback.MenuItem {
	case types.MenuItemTest1:
		s.log.Debug("chooseMenuItem - MenuItemTest1")
	case types.MenuItemTest2:
		s.log.Debug("chooseMenuItem - MenuItemTest2")
	case types.MenuItemTest3:
		s.log.Debug("chooseMenuItem - MenuItemTest3")
	}
	return nil
}
