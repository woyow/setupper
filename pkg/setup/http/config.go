package http

import (
	"github.com/woyow/setupper/pkg/setup/http/client"
)

type Config struct {
	Client client.Config `yaml:"client"`
}
