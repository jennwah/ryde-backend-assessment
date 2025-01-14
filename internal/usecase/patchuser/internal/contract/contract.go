package contract

import (
	"context"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
)

type Repository interface {
	Patch(ctx context.Context, user model.PatchUser) error
}
