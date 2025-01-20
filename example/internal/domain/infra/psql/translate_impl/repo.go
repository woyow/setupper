package translate_impl

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	domainLoggingKey   = "domain"
	domainLoggingValue = "your-tg"
	queryLoggingKey    = "query"
	argsLoggingKey     = "args"
)

type Repo struct {
	db  *pgxpool.Pool
	qb  *squirrel.StatementBuilderType
	log *logrus.Logger
}

func NewRepo(db *pgxpool.Pool, qb *squirrel.StatementBuilderType, log *logrus.Logger) *Repo {
	return &Repo{
		db:  db,
		qb:  qb,
		log: log,
	}
}
