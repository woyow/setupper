package translate_impl

import (
	"context"

	"github.com/woyow/setupper/pkg/translate"

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

func (r *Repo) GetLocalizationMap(ctx context.Context, dbName string) (out translate.LocalizationMap, err error) {
	query := "SELECT jsonb_object_agg(t.key, t.name) FROM " + dbName + " AS t"

	if err := r.db.QueryRow(ctx, query).Scan(&out); err != nil {
		r.log.WithFields(logrus.Fields{
			domainLoggingKey: domainLoggingValue,
			queryLoggingKey:  query,
			argsLoggingKey:   nil,
		}).Error("psql: GetLocalizationMap query error: ", err)
		return nil, err
	}

	return out, nil
}
