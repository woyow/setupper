package migrate

import (
	"errors"
	"time"

	setupPsql "github.com/woyow/setupper/pkg/setup/psql"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

type Migrate struct {
	driver          database.Driver
	log             *logrus.Logger
	source          string
	dbname          string
	dburl           string
	attempts        int
	attemptsTimeout time.Duration
	isEnabled       bool
}

const (
	setupLoggingKey   = "setup"
	setupLoggingValue = "postgresql migrate"
)

func NewMigrate(setupPsql *setupPsql.Psql, log *logrus.Logger) (*Migrate, error) {
	attempts, attemptsTimeout := setupPsql.GetMigrationAttempts()

	return &Migrate{
		source:          setupPsql.GetMigrationSource(),
		dbname:          setupPsql.GetDatabaseName(),
		dburl:           setupPsql.GetDatabaseURL(),
		attempts:        attempts,
		attemptsTimeout: attemptsTimeout,
		log:             log,
		isEnabled:       setupPsql.IsMigrationEnabled(),
	}, nil
}

func (m *Migrate) Run() error {
	if !m.isEnabled {
		m.log.WithField(setupLoggingKey, setupLoggingValue).
			Info("Postgresql migration disabled for db: ", m.dbname)
		return nil
	}
	var (
		p   = pgx.Postgres{}
		err error
	)

	for m.attempts > 0 {
		m.driver, err = p.Open(m.dburl)
		if err == nil {
			break
		}

		<-time.After(m.attemptsTimeout)
		m.attempts--
	}
	if err != nil {
		m.log.WithField(setupLoggingKey, setupLoggingValue).
			Error("Run - p.Open error: ", err.Error())
		return err
	}

	mm, err := migrate.NewWithDatabaseInstance(m.source, m.dbname, m.driver)
	if err != nil {
		m.log.WithField(setupLoggingKey, setupLoggingValue).
			Error("NewMigrate - migrate.NewWithDatabaseInstance error: ", err.Error())
		return err
	}

	defer mm.Close()

	if !m.isEnabled {
		m.log.WithField(setupLoggingKey, setupLoggingValue).
			Info("Run - migration is disabled")
		return nil
	}

	if err := mm.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.log.WithField(setupLoggingKey, setupLoggingValue).
				Info("Run - no change")

			return nil
		}

		m.log.WithField(setupLoggingKey, setupLoggingValue).
			Fatal("Run - mm.Up error: ", err.Error())

		return err
	}

	return nil
}
