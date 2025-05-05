package inspector

import (
	"context"
	"github.com/florianl/go-nfqueue"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"obfs-detector/internal/db"
	"obfs-detector/internal/model"
	"obfs-detector/pkg/detector"
	"sync"
)

type inspector struct {
	logger           *zerolog.Logger
	postgres         db.Postgres
	blockedAddresses *sync.Map
	queue            *nfqueue.Nfqueue
	detectors        []detector.Detector
}

func New(logger *zerolog.Logger, postgres db.Postgres, detectors []detector.Detector, queue *nfqueue.Nfqueue) *inspector {
	return &inspector{
		logger:    logger,
		postgres:  postgres,
		detectors: detectors,
		queue:     queue,
	}
}

func (i *inspector) Start(ctx context.Context) error {
	addressesSlice, err := i.postgres.SelectBlockedIPAddresses(ctx)
	if err != nil {
		return errors.Wrap(err, "select blocked ip address from postgres fail")
	}

	addressesMap := sync.Map{}
	for _, address := range addressesSlice {
		addressesMap.Store(address.Address, address)
	}

	// Register your function to listen on nflqueue queue 100
	err = i.queue.RegisterWithErrorFunc(ctx, i.handler, func(e error) int {
		i.logger.Warn().Err(e).Msg("set accept verdict for clear packet fail")
		return -1
	})
	if err != nil {
		return errors.Wrap(err, "attaches callback function and error function to nfqueue fail")
	}

	return nil
}

func (i *inspector) Stop(ctx context.Context) error {
	addresses := make([]model.BlockedIPAddress, 0)
	i.blockedAddresses.Range(func(key, value interface{}) bool {
		address, ok := value.(model.BlockedIPAddress)
		if ok {
			return false
		}

		addresses = append(addresses, address)
		return true
	})

	if err := i.postgres.UpsertBlockedIPAddresses(ctx, addresses); err != nil {
		return errors.Wrap(err, "upsert blocked ip addresses fail")
	}

	return nil
}
