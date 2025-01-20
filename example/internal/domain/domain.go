package domain

import (
	telegramDelivery "github.com/woyow/setupper/example/internal/domain/delivery/telegram"
	"github.com/woyow/setupper/example/internal/domain/domain/service"
	"github.com/woyow/setupper/example/internal/domain/infra"

	"github.com/woyow/setupper/pkg/telegram"
	"github.com/woyow/setupper/pkg/translate"
	"time"

	setupEchotron "github.com/woyow/setupper/pkg/setup/echotron"
	setupPsql "github.com/woyow/setupper/pkg/setup/psql"

	"github.com/sirupsen/logrus"
)

const (
	domainLoggingKey   = "domain"
	domainLoggingValue = "your-tg"
)

type Setup struct {
	Psql     *setupPsql.Psql
	Echotron *setupEchotron.Echotron
}

func NewDomain(setup *Setup, stop <-chan struct{}, log *logrus.Logger) {

	infr := infra.NewInfra(setup.Psql, setup.Echotron, log)

	t := translate.NewTranslate(infr.Psql.TranslateImpl, translate.Config{
		UpdateTimeout:      1 * time.Minute,
		DBName:             "your_tg_translates",
		DefaultLanguage:    "ru",
		AvailableLanguages: []string{"ru", "en"},
	}, stop, log)

	svc := service.NewService(infr, t, log)

	tg := telegram.NewTelegram(
		setup.Echotron, 
		infr.Psql.StateImpl, 
		log, 
		domainLoggingValue,
	)

	telegramDelivery.InitHandler(tg, svc, log)

	log.WithFields(logrus.Fields{
		domainLoggingKey: domainLoggingValue,
	}).Info("domain has been initialized")
}
