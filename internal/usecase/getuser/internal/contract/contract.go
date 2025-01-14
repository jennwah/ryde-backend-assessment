package contract

import (
	"context"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
)

type Repository interface {
	Get(ctx context.Context, id string) (model.User, error)
}
