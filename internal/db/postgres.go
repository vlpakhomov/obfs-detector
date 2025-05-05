package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"obfs-detector/internal/model"
)

type postgres struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *postgres {
	return &postgres{
		db: db,
	}
}

func (p *postgres) SelectBlockedIPAddresses(ctx context.Context) ([]model.BlockedIPAddress, error) {
	addresses := make([]model.BlockedIPAddress, 0)

	rows, err := p.db.QueryxContext(ctx, SelectBlockedIPAddressesSQL)
	if err != nil {
		return nil, errors.Wrap(err, "select blocked ip addresses fail")
	}
	// TODO: handle error from rows.Close()
	defer rows.Close()

	for rows.Next() {
		address := model.BlockedIPAddress{}
		if err = rows.StructScan(&address); err != nil {
			return nil, errors.Wrap(err, "scan row to blocked ip address fail")
		}

		addresses = append(addresses, address)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(err, "iteration over db rows fail")
	}

	return addresses, nil
}

func (p *postgres) UpsertBlockedIPAddresses(ctx context.Context, addresses []model.BlockedIPAddress) error {
	_, err := p.db.NamedExecContext(ctx, UpsertBlockedIPAddressesSQL, addresses)
	if err != nil {
		return errors.Wrap(err, "upsert blocked ip addresses fail")
	}

	return nil
}
