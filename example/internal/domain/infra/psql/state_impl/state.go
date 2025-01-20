package state_impl

import (
	"context"

	"github.com/woyow/setupper/pkg/telegram"
	"github.com/sirupsen/logrus"
)

func (r *Repo) IsBanned(ctx context.Context, dto telegram.IsBannedDTO) (out bool) {
	query := "SELECT your_tg_check_ban($1::varchar)"

	if err := r.db.QueryRow(ctx, query, dto.ChatID).Scan(&out); err != nil {
		r.log.WithFields(logrus.Fields{
			domainLoggingKey: domainLoggingValue,
			queryLoggingKey:  query,
			argsLoggingKey: []interface{}{
				dto.ChatID,
			},
		}).Error("psql: IsBanned query error: ", err)
		return out
	}

	return out
}

func (r *Repo) CreateCurrentState(ctx context.Context, dto telegram.CreateCurrentStateDTO) error {
	query := "INSERT INTO your_tg_chat_states(chat_id, state) VALUES ($1, $2)"

	if _, err := r.db.Exec(ctx, query, dto.ChatID, dto.State); err != nil {
		r.log.WithFields(logrus.Fields{
			domainLoggingKey: domainLoggingValue,
			queryLoggingKey:  query,
			argsLoggingKey: []interface{}{
				dto.ChatID,
				dto.State,
			},
		}).Error("psql: CreateCurrentState query error: ", err)
		return err
	}

	return nil
}

func (r *Repo) SetCurrentState(ctx context.Context, dto telegram.SetCurrentStateDTO) error {
	query := "INSERT INTO your_tg_chat_states(chat_id, state) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET state = excluded.state"

	if _, err := r.db.Exec(ctx, query, dto.ChatID, dto.State); err != nil {
		r.log.WithFields(logrus.Fields{
			domainLoggingKey: domainLoggingValue,
			queryLoggingKey:  query,
			argsLoggingKey: []interface{}{
				dto.ChatID,
				dto.State,
			},
		}).Error("psql: SetCurrentState query error: ", err)
		return err
	}

	return nil
}

const (
	noRowsInResultSetError = "no rows in result set"
)

func (r *Repo) GetCurrentState(ctx context.Context, dto telegram.GetCurrentStateDTO) (out telegram.GetCurrentStateOut, err error) {
	query := "SELECT state FROM your_tg_chat_states WHERE chat_id = $1"

	if err := r.db.QueryRow(ctx, query, dto.ChatID).Scan(&out.State); err != nil {
		if err.Error() == noRowsInResultSetError {
			return out, telegram.ErrChatCurrentStateNotExists
		}
		r.log.WithFields(logrus.Fields{
			domainLoggingKey: domainLoggingValue,
			queryLoggingKey:  query,
			argsLoggingKey: []interface{}{
				dto.ChatID,
			},
		}).Error("psql: GetCurrentState query error: ", err)
		return out, err
	}

	return out, nil
}
