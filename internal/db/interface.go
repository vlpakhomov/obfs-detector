package db

import (
	"context"
	"obfs-detector/internal/model"
)

type Postgres interface {
	SelectBlockedIPAddresses(ctx context.Context) ([]model.BlockedIPAddress, error)
	UpsertBlockedIPAddresses(ctx context.Context, addresses []model.BlockedIPAddress) error
}
