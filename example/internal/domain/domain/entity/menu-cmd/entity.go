package menu_cmd

import (
	"github.com/woyow/setupper/pkg/telegram"
)

type MenuCommandServiceDTO struct {
	telegram.HandleCommand
}

type MenuCommandCallbackServiceDTO struct {
	telegram.HandleCallback
	CallbackData string
}

type MenuCommandCallbackOut struct {
	WaitCallback bool
}
