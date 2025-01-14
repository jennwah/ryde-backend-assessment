package contract

import (
	"context"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (string, error)
}
