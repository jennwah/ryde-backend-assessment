package contract

import (
	"context"
)

type Repository interface {
	Delete(ctx context.Context, id string) error
}
