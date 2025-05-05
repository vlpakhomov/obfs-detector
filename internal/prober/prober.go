package prober

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net"
	"obfs-detector/internal/db"
	"obfs-detector/pkg/null"
	"obfs-detector/pkg/obfs4/transports/obfs3"
	"time"
)

type prober struct {
	logger *zerolog.Logger
	db     db.Postgres
	client obfs3.Transport
}

func New(logger *zerolog.Logger, db db.Postgres, client obfs3.Transport) Prober {
	return &prober{
		logger: logger,
		db:     db,
		client: client,
	}
}

func (p *prober) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			addresses, err := p.db.SelectBlockedIPAddresses(ctx)
			if err != nil {
				p.logger.Warn().Msg("select blocked ip addresses fail")
				continue
			}

			for _, address := range addresses {
				detected, err := p.probe(address.Address)
				if err != nil {
					p.logger.Warn().Msg("probe blocked ip addresses fail")
					continue
				}

				if *detected.ValuePtr() {
					p.logger.Info().Msg("obfs3 detected")
					continue
				}
			}
		}
	}
}

func (p *prober) probe(address string) (null.Null[bool], error) {
	factory, err := p.client.ClientFactory("")
	if err != nil {
		p.logger.Info().Msg("create new obfs3 client factory instance fail")

		return null.NewExplicit(false, false),
			errors.Wrap(err, "create new obfs3 client factory instance fail")
	}

	_, err = factory.Dial("tcp", address, net.Dial, "")
	if err != nil {
		return null.New(true), nil
	}

	return null.New(false), nil
}
