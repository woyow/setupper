package telegram

import "errors"

var (
	ErrChatCurrentStateNotExists     = errors.New("chat current state not exists")
	ErrChatCurrentStateAlreadyExists = errors.New("chat current state already exists")
)
