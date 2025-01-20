package menu_cmd

import (
	types "github.com/woyow/setupper/example/internal/domain/domain/types/menu-cmd"
)

type CallbackOut struct {
	Type types.CallbackType `json:"t"`
}

type ChooseMenuCallbackOut struct {
	CallbackOut
}

type ChooseMenuItemCallbackOut struct {
	CallbackOut
	MenuItem types.MenuItem `json:"i"`
}
