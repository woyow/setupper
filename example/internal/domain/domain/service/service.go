package service

import (
	menuCmd "github.com/woyow/setupper/example/internal/domain/domain/service/menu-cmd"

	"github.com/woyow/setupper/example/internal/domain/infra"
	"github.com/woyow/setupper/pkg/translate"

	"github.com/sirupsen/logrus"
)

type Service struct {
	MenuCmd *menuCmd.Service
}

func NewService(infra *infra.Infra, translate *translate.Translate, log *logrus.Logger) *Service {
	return &Service{
		MenuCmd: menuCmd.NewService(
			nil,
			infra.TgApi,
			translate,
			log,
		),
	}
}