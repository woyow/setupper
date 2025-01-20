package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/woyow/setupper/example/internal/config"
	"github.com/woyow/setupper/example/internal/domain"

	setupPsql "github.com/woyow/setupper/pkg/setup/psql"
	setupLogger "github.com/woyow/setupper/pkg/setup/logger"
	setupEchotron "github.com/woyow/setupper/pkg/setup/echotron"

	psqlMigrate "github.com/woyow/setupper/pkg/setup/psql/migrate"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type setup struct {
	domain *domain.Setup
	psql   *setupPsql.Psql
}

type migrate struct {
	psql *psqlMigrate.Migrate
}

type app struct {
	log      *logrus.Logger
	cfg      *config.Config
	errGroup *errgroup.Group
	setup    setup
	migrate  migrate
	stopCh   chan os.Signal
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewApp() *app {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background()) // Base app context

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger := setupLogger.NewLogger(&cfg.Logger)

	errGroup, ctx := errgroup.WithContext(ctx)

	psql, err := setupPsql.NewPsql(ctx, &cfg.Psql, stopCh, logger)
	if err != nil {
		panic(err)
	}

	psqlMigration, err := psqlMigrate.NewMigrate(psql, logger)
	if err != nil {
		panic(err)
	}

	yourTgEchotron, err := setupEchotron.NewEchotron(&cfg.YourTg.Echotron, logger)
	if err != nil {
		panic(err)
	}

	return &app{
		log:      logger,
		cfg:      cfg,
		errGroup: errGroup,
		stopCh:   stopCh,
		ctx:      ctx,
		cancelFn: cancel,
		setup: setup{
			domain: &domain.Setup{
				Echotron:      yourTgEchotron,
				Psql:          psql,
			},
			psql:  psql,
		},
		migrate: migrate{
			psql: psqlMigration,
		},
	}
}

func (a *app) Run() error {
	// Run migrations
	{
		if err := a.migrate.psql.Run(); err != nil {
			return err
		}
	}

	// Initialize domains
	{
		domain.NewDomain(a.setup.domain, a.stopCh, a.log)
	}

	// Handle stop program
	a.errGroup.Go(func() error {
		a.log.Infof("Got %s signal. Aborting...\n", <-a.stopCh)
		a.cancelFn()
		close(a.stopCh)
		return nil
	})

	// Wait error from group of goroutines
	if err := a.errGroup.Wait(); err != nil {
		a.log.Error("app: Run - g.Wait error: ", err)
		return err
	}

	return nil
}
