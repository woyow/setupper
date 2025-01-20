package common

type GetFileAPIDTO struct {
	FileID string
}

type GetFileAPIOut struct {
	FilePath string
}

type Options struct {
	ParseMode      *string
	RemoveKeyboard *bool
	DisablePreview *bool
}

type SendMessageAPIDTO struct {
	Text    string
	ChatID  int64
	Options *Options
}

type InlineButton struct {
	Text         string
	CallbackData string
	URL          string
}

type SendMessageWithInlineKeyboardAPIDTO struct {
	Buttons [][]InlineButton
	Options *Options
	Text    string
	ChatID  int64
}

type EditMessageWithInlineKeyboardAPIDTO struct {
	Buttons   [][]InlineButton
	Options   *Options
	ChatID    int64
	MessageID int
}

type EditMessageTextAPIDTO struct {
	Buttons   [][]InlineButton
	Options   *Options
	Text      string
	ChatID    int64
	MessageID int
}

type KeyboardButton struct {
	Text           string
	RequestContact bool
}

type SendMessageWithKeyboardAPIDTO struct {
	Buttons         [][]KeyboardButton
	Text            string
	ChatID          int64
	OneTimeKeyboard bool
}

type DeleteMessageAPIDTO struct {
	ChatID    int64
	MessageID int
}

type DeleteMessagesAPIDTO struct {
	ChatID     int64
	MessageIDs []int
}
