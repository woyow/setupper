package telegram

type HandleCommand struct {
	Lang   string
	ChatID int64
}

type HandleMessage struct {
	Lang      string
	ChatID    int64
	MessageID int
}

type HandleCallback struct {
	Lang      string
	ChatID    int64
	MessageID int
}

type IsBannedDTO struct {
	ChatID int64
}

type CreateCurrentStateDTO struct {
	State  string
	ChatID int64
}

type SetCurrentStateDTO struct {
	State  string
	ChatID int64
}

type GetCurrentStateDTO struct {
	ChatID int64
}

type GetCurrentStateOut struct {
	State string
}
