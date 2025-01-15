package contract

import (
	"context"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
)

type Repository interface {
	Get(ctx context.Context, userID string, radiusMeter float64) (model.Friends, error)
}
