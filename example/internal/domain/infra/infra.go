package infra

import (
	"github.com/woyow/setupper/example/internal/domain/infra/psql"
	tgapiImpl "github.com/woyow/setupper/example/internal/domain/infra/tgapi_impl"

	setupEchotron "github.com/woyow/setupper/pkg/setup/echotron"
	setupPsql "github.com/woyow/setupper/pkg/setup/psql"

	"github.com/sirupsen/logrus"
)

type Infra struct {
	Psql  *psql.Psql
	TgApi *tgapiImpl.API
}

func NewInfra(
	setupPsql *setupPsql.Psql,
	setupEchotron *setupEchotron.Echotron,
	log *logrus.Logger) *Infra {
	return &Infra{
		Psql:  psql.NewPsql(setupPsql, log),
		TgApi: tgapiImpl.NewAPI(setupEchotron.GetAPI(), log),
	}
}
