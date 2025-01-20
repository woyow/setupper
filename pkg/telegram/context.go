package telegram

import "github.com/NicoNex/echotron/v3"

type Context echotron.Update

func (c *Context) ChatID() int64 {
	return echotron.Update(*c).ChatID()
}
