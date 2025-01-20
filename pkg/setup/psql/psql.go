package psql

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	setupLoggingKey   = "setup"
	setupLoggingValue = "psql"

	proto = "postgres"
)

var (
	errParseConfig = errors.New("pgxpool parse config error")
)

type Psql struct {
	queryBuilder squirrel.StatementBuilderType // Query builder
	pool         *pgxpool.Pool                 // Pool of connections
	log          *logrus.Logger
	migration    migration
	databaseURL  string
	databaseName string
	stop         <-chan struct{}
}

type migration struct {
	attemptTimeout time.Duration
	source         string
	attempts       int
	isEnabled      bool
}

func NewPsql(ctx context.Context, cfg *Config, stop <-chan struct{}, log *logrus.Logger) (*Psql, error) {
	databaseURL := getDatabaseURL(cfg)

	queryBuilder := getQueryBuilder()

	pool, err := getPool(ctx, cfg, log)
	if err != nil {
		return nil, err
	}

	p := &Psql{
		queryBuilder: queryBuilder,
		pool:         pool,
		log:          log,
		databaseName: os.Getenv(cfg.DBNameEnvKey),
		databaseURL:  databaseURL,
		migration: migration{
			source:         cfg.Migration.Source,
			attemptTimeout: time.Duration(cfg.Migration.AttemptsTimeout) * time.Second,
			attempts:       cfg.Migration.Attempts,
			isEnabled:      cfg.Migration.Enable,
		},
		stop: stop,
	}

	go func() {
		if err := p.shutdown(); err != nil {
			p.log.WithFields(logrus.Fields{
				setupLoggingKey: setupLoggingValue,
			}).Error("p.shutdown error: ", err)
		}
	}()

	return p, nil
}

// getQueryBuilder - Returns squirrel query builder
func getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// getPool - Returns pool of connections to postgresql database
func getPool(ctx context.Context, cfg *Config, log *logrus.Logger) (*pgxpool.Pool, error) {
	databaseURL := getDatabaseURL(cfg)

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.WithField(setupLoggingKey, setupLoggingValue).
			Error("getPool - pgxpool.ParseConfig error: ", err.Error())
		return nil, errParseConfig
	}

	{
		poolConfig.MaxConns = int32(cfg.Pool.MaxPoolSize)
		poolConfig.MinConns = int32(cfg.Pool.MinPoolSize)

		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.WithField(setupLoggingKey, setupLoggingValue).
			Error("getPool - pgxpool.NewWithConfig error: ", err.Error())
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.WithField(setupLoggingKey, setupLoggingValue).
			Error("getPool - pool.Ping error: ", err.Error())
		return nil, err
	}

	return pool, nil
}

func getDatabaseURL(cfg *Config) string {
	return fmt.Sprintf(
		"%s/%s%s",
		os.Getenv(cfg.URLEnvKey),
		os.Getenv(cfg.DBNameEnvKey),
		func() string {
			parameters := os.Getenv(cfg.URLParametersEnvKey)
			if parameters == "" {
				return ""
			} else {
				return "?" + parameters
			}
		}(),
	)
}

func (p *Psql) IsMigrationEnabled() bool {
	return p.migration.isEnabled
}

func (p *Psql) GetDatabaseName() string {
	return p.databaseName
}

func (p *Psql) GetDatabaseURL(options ...string) string {
	if options != nil && len(options) > 0 {
		return fmt.Sprintf("%s&%s", p.databaseURL, strings.Join(options, "&"))
	} else {
		return fmt.Sprintf("%s", p.databaseURL)
	}
}

func (p *Psql) GetMigrationAttempts() (int, time.Duration) {
	return p.migration.attempts, p.migration.attemptTimeout
}

func (p *Psql) GetMigrationSource() string {
	return p.migration.source
}

func (p *Psql) GetQueryBuilder() *squirrel.StatementBuilderType {
	return &p.queryBuilder
}

func (p *Psql) GetPool() *pgxpool.Pool {
	return p.pool
}

func (p *Psql) shutdown() error {
	select {
	case <-p.stop:
		p.log.WithField(setupLoggingKey, setupLoggingValue).
			Info("Shutdown - close postgresql pool")

		p.pool.Close()
	}
	return nil
}
