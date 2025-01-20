package translate

import (
	"time"
)

type Config struct {
	UpdateTimeout      time.Duration
	DBName             string
	DefaultLanguage    string
	AvailableLanguages []string
}