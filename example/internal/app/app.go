package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/woyow/setupper/example/internal/config"
	"github.com/woyow/setupper/example/internal/domain"

	setupEchotron "github.com/woyow/setupper/pkg/setup/echotron"
	setupLogger "github.com/woyow/setupper/pkg/setup/logger"
	setupPsql "github.com/woyow/setupper/pkg/setup/psql"

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
	sigCh    chan os.Signal
	stopCh   chan struct{}
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewApp() *app {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	stopCh := make(chan struct{}, 1)

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
		sigCh:    sigCh,
		stopCh:   stopCh,
		ctx:      ctx,
		cancelFn: cancel,
		setup: setup{
			domain: &domain.Setup{
				Echotron: yourTgEchotron,
				Psql:     psql,
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

	// Initialize domain
	{
		domain.NewDomain(a.setup.domain, a.stopCh, a.log)
	}

	// Handle stop program
	a.errGroup.Go(func() error {
		a.log.Infof("Got %s signal. Aborting...\n", <-a.sigCh)
		a.cancelFn()
		close(a.sigCh)
		close(a.stopCh)
		<-time.After(1 * time.Second)
		return nil
	})

	// Wait error from group of goroutines
	if err := a.errGroup.Wait(); err != nil {
		a.log.Error("app: Run - g.Wait error: ", err)
		return err
	}

	return nil
}
