package prober

import (
	"context"
)

type Prober interface {
	Start(ctx context.Context) error
}
