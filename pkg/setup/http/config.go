package http

import (
	"shared-lib/pkg/setup/http/client"
)

type Config struct {
	Client client.Config `yaml:"client"`
}
