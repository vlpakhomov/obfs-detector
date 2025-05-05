package inspector

import "context"

type Inspector interface {
	Start(ctx context.Context) error
}
