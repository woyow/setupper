package translate_impl

import (
	"context"

	"github.com/woyow/setupper/pkg/translate"

	"github.com/sirupsen/logrus"
)

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