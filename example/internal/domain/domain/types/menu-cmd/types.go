package menu_cmd

type CallbackType int16

const (
	ChooseMenu CallbackType = iota + 1
	ChooseMenuItem
)

type MenuItem int16

const (
	MenuItemTest1 MenuItem = iota + 1
	MenuItemTest2
	MenuItemTest3
)