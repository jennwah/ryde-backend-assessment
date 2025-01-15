package contract

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, userID, friendID string) error
}
