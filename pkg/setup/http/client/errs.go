package client

import "errors"

var (
	ErrRedirectNotAllowed = errors.New("redirect not allowed")
)
