package app

import (
	"context"
	"github.com/florianl/go-nfqueue"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"obfs-detector/internal/db"
	"obfs-detector/internal/inspector"
	"obfs-detector/internal/prober"
	"obfs-detector/pkg/detector"
	"obfs-detector/pkg/obfs4/transports/obfs3"
	"os/signal"
	"sync"
	"syscall"
)

type app struct {
	logger    *zerolog.Logger
	db        *sqlx.DB
	queue     *nfqueue.Nfqueue
	inspector inspector.Inspector
}

func New(logger *zerolog.Logger, db *sqlx.DB, queue *nfqueue.Nfqueue) *app {
	return &app{
		logger: logger,
		db:     db,
		queue:  queue,
	}
}

func (a *app) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	postgres := db.New(a.db)

	insp := inspector.New(a.logger, postgres, []detector.Detector{&detector.Obfs2, &detector.Obfs3}, a.queue)
	a.inspector = insp
	if err := a.inspector.Start(ctx); err != nil {
		a.logger.Fatal().Err(err).Msg("start obfs inspector fail")
	}

	wg := sync.WaitGroup{}

	pbr := prober.New(a.logger, postgres, obfs3.Transport{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		pbr.Start(ctx)
	}()

	wg.Wait()
}
