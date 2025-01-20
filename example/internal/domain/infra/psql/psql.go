package psql

import (
	stateImpl "github.com/woyow/setupper/example/internal/domain/infra/psql/state_impl"
	translateImpl "github.com/woyow/setupper/example/internal/domain/infra/psql/translate_impl"

	setupPsql "github.com/woyow/setupper/pkg/setup/psql"

	"github.com/sirupsen/logrus"
)

type Psql struct {
	StateImpl     *stateImpl.Repo
	TranslateImpl *translateImpl.Repo
}

func NewPsql(setupPsql *setupPsql.Psql, log *logrus.Logger) *Psql {
	return &Psql{
		StateImpl:     stateImpl.NewRepo(setupPsql.GetPool(), setupPsql.GetQueryBuilder(), log),
		TranslateImpl: translateImpl.NewRepo(setupPsql.GetPool(), setupPsql.GetQueryBuilder(), log),
	}
}
