package tgapi

import (
	commonEntity "github.com/woyow/setupper/example/internal/domain/domain/entity/common"
	
	"github.com/woyow/setupper/example/internal/domain/errs"

	"github.com/NicoNex/echotron/v3"
	"github.com/sirupsen/logrus"
)

func (a *API) GetFile(dto commonEntity.GetFileAPIDTO) (commonEntity.GetFileAPIOut, error) {
	resp, err := a.api.GetFile(dto.FileID)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
			fileIDLoggingKey: dto.FileID,
		}).Error("GetFile error: ", err)
		return commonEntity.GetFileAPIOut{}, err
	}

	out := commonEntity.GetFileAPIOut{
		FilePath: resp.Result.FilePath,
	}

	return out, nil
}


func (a *API) getInlineKeyboard(buttons [][]commonEntity.InlineButton) [][]echotron.InlineKeyboardButton {
	inlineKeyboard := make([][]echotron.InlineKeyboardButton, 0, len(buttons))

	for i := range buttons {
		inlineButtons := make([]echotron.InlineKeyboardButton, 0, len(buttons[i]))

		for j := range buttons[i] {
			inlineButtons = append(inlineButtons, echotron.InlineKeyboardButton{
				CallbackGame:                 nil,
				WebApp:                       nil,
				LoginURL:                     nil,
				SwitchInlineQueryChosenChat:  nil,
				Text:                         buttons[i][j].Text,
				CallbackData:                 buttons[i][j].CallbackData,
				SwitchInlineQuery:            "",
				SwitchInlineQueryCurrentChat: "",
				URL:                          buttons[i][j].URL,
				Pay:                          false,
			})
		}

		inlineKeyboard = append(inlineKeyboard, inlineButtons)
	}

	return inlineKeyboard
}

func (a *API) SendMessageWithInlineKeyboard(dto commonEntity.SendMessageWithInlineKeyboardAPIDTO) error {
	inlineKeyboard := a.getInlineKeyboard(dto.Buttons)

	options := &echotron.MessageOptions{
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
		ParseMode: func() echotron.ParseMode {
			if dto.Options != nil && dto.Options.ParseMode != nil {
				return echotron.ParseMode(*dto.Options.ParseMode)
			} else {
				return echotron.Markdown
			}
		}(),
		LinkPreviewOptions: echotron.LinkPreviewOptions{
			URL: "",
			IsDisabled: func() bool {
				if dto.Options != nil && dto.Options.DisablePreview != nil {
					return *dto.Options.DisablePreview
				} else {
					return false
				}
			}(),
			PreferSmallMedia: false,
			PreferLargeMedia: false,
			ShowAboveText:    false,
		},
	}

	if dto.Options != nil {
		if dto.Options.RemoveKeyboard != nil {
			options.ReplyMarkup = echotron.ReplyKeyboardRemove{
				RemoveKeyboard: *dto.Options.RemoveKeyboard,
				Selective:      false,
			}
		}
	}

	resp, err := a.api.SendMessage(dto.Text, dto.ChatID, options)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
			textLoggingKey:   dto.Text,
		}).Error("SendMessageWithInlineKeyboard - s.api.SendMessage error", err)

		return err
	}

	if !resp.Ok {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
			textLoggingKey:   dto.Text,
		}).Error("SendMessageWithInlineKeyboard - status code: ", resp.ErrorCode)

		return errs.ErrStatusCodeUnsuccessful
	}

	return nil
}

const sameContentError = "API error: 400 Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message"

func (a *API) EditMessageWithInlineKeyboard(dto commonEntity.EditMessageWithInlineKeyboardAPIDTO) error {
	inlineKeyboard := a.getInlineKeyboard(dto.Buttons)

	msgIDOpt := echotron.NewMessageID(dto.ChatID, dto.MessageID)

	resp, err := a.api.EditMessageReplyMarkup(msgIDOpt, &echotron.MessageReplyMarkupOptions{
		BusinessConnectionID: "",
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	if err != nil {
		if err.Error() == sameContentError {
			return nil
		}
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
		}).Error("EditMessageWithInlineKeyboard - s.api.SendMessage error: ", err)

		return err
	}

	if !resp.Ok {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
		}).Error("EditMessageWithInlineKeyboard - status code: ", resp.ErrorCode)

		return errs.ErrStatusCodeUnsuccessful
	}

	return nil
}

func (a *API) EditMessageText(dto commonEntity.EditMessageTextAPIDTO) error {

	inlineKeyboard := a.getInlineKeyboard(dto.Buttons)

	msgIDOpt := echotron.NewMessageID(dto.ChatID, dto.MessageID)

	options := &echotron.MessageTextOptions{
		ParseMode: func() echotron.ParseMode {
			if dto.Options != nil && dto.Options.ParseMode != nil {
				return echotron.ParseMode(*dto.Options.ParseMode)
			} else {
				return "markdown"
			}
		}(),
		LinkPreviewOptions: echotron.LinkPreviewOptions{
			URL: "",
			IsDisabled: func() bool {
				if dto.Options != nil && dto.Options.DisablePreview != nil {
					return *dto.Options.DisablePreview
				} else {
					return false
				}
			}(),
			PreferSmallMedia: false,
			PreferLargeMedia: false,
			ShowAboveText:    false,
		},
		Entities: nil,
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	}

	resp, err := a.api.EditMessageText(dto.Text, msgIDOpt, options)
	if err != nil {
		if err.Error() == sameContentError {
			return nil
		}
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
		}).Error("EditMessageText - s.api.SendMessage error: ", err)

		return err
	}

	if !resp.Ok {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
		}).Error("EditMessageText - status code: ", resp.ErrorCode)

		return errs.ErrStatusCodeUnsuccessful
	}

	return nil
}

func (a *API) getKeyboard(buttons [][]commonEntity.KeyboardButton) [][]echotron.KeyboardButton {
	keyboard := make([][]echotron.KeyboardButton, 0, len(buttons))

	for i := range buttons {
		keyboardButtons := make([]echotron.KeyboardButton, 0, len(buttons[i]))

		for j := range buttons[i] {
			keyboardButtons = append(keyboardButtons, echotron.KeyboardButton{
				RequestPoll:     nil,
				WebApp:          nil,
				RequestUsers:    nil,
				RequestChat:     nil,
				Text:            buttons[i][j].Text,
				RequestContact:  buttons[i][j].RequestContact,
				RequestLocation: false,
			})
		}

		keyboard = append(keyboard, keyboardButtons)
	}

	return keyboard
}

func (a *API) SendMessageWithKeyboard(dto commonEntity.SendMessageWithKeyboardAPIDTO) error {
	keyboard := a.getKeyboard(dto.Buttons)

	if _, err := a.api.SendMessage(dto.Text, dto.ChatID, &echotron.MessageOptions{
		ReplyMarkup: echotron.ReplyKeyboardMarkup{
			InputFieldPlaceholder: "",
			Keyboard:              keyboard,
			IsPersistent:          false,
			ResizeKeyboard:        false,
			OneTimeKeyboard:       dto.OneTimeKeyboard,
			Selective:             true,
		},
		BusinessConnectionID: "",
		MessageEffectID:      "",
		ParseMode:            "markdown",
		LinkPreviewOptions:   echotron.LinkPreviewOptions{},
		Entities:             nil,
		ReplyParameters:      echotron.ReplyParameters{},
		MessageThreadID:      0,
		DisableNotification:  false,
		ProtectContent:       false,
	}); err != nil {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
		}).Error("SendMessageWithKeyboard - s.api.SendMessage error: ", err)
		return err
	}

	return nil
}

func (a *API) DeleteMessage(dto commonEntity.DeleteMessageAPIDTO) error {
	_, err := a.api.DeleteMessage(dto.ChatID, dto.MessageID)
	if err != nil {
		a.log.WithFields(logrus.Fields{}).Error("DeleteMessage - s.api.DeleteMessage error: ", err)
		return err
	}

	return nil
}

func (a *API) DeleteMessages(dto commonEntity.DeleteMessagesAPIDTO) error {
	_, err := a.api.DeleteMessages(dto.ChatID, dto.MessageIDs)
	if err != nil {
		a.log.WithFields(logrus.Fields{}).Error("DeleteMessages - s.api.DeleteMessages error: ", err)
		return err
	}

	return nil
}

func (a *API) SendMessage(dto commonEntity.SendMessageAPIDTO) error {
	options := &echotron.MessageOptions{
		ReplyMarkup:          nil,
		BusinessConnectionID: "",
		MessageEffectID:      "",
		ParseMode: func() echotron.ParseMode {
			if dto.Options != nil && dto.Options.ParseMode != nil {
				return echotron.ParseMode(*dto.Options.ParseMode)
			} else {
				return "markdown"
			}
		}(),
		LinkPreviewOptions:  echotron.LinkPreviewOptions{},
		Entities:            nil,
		ReplyParameters:     echotron.ReplyParameters{},
		MessageThreadID:     0,
		DisableNotification: false,
		ProtectContent:      false,
	}

	if dto.Options != nil {
		if dto.Options.RemoveKeyboard != nil {
			options.ReplyMarkup = echotron.ReplyKeyboardRemove{
				RemoveKeyboard: *dto.Options.RemoveKeyboard,
				Selective:      false,
			}
		}
	}

	resp, err := a.api.SendMessage(dto.Text, dto.ChatID, options)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
			textLoggingKey:   dto.Text,
		}).Error("SendMessage - s.api.SendMessage error: ", err)

		return err
	}

	if !resp.Ok {
		a.log.WithFields(logrus.Fields{
			chatIDLoggingKey: dto.ChatID,
			infraLoggingKey:  infraLoggingValue,
			tgBotNameLoggingKey: tgBotNameLoggingValue,
			textLoggingKey:   dto.Text,
		}).Error("SendMessage - status code: ", resp.ErrorCode)

		return errs.ErrStatusCodeUnsuccessful
	}

	return nil
}
