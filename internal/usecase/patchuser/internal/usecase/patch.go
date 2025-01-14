package usecase

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
)

func (u UseCase) Patch(ctx context.Context, user model.PatchUser) error {
	err := u.repo.Patch(ctx, user)
	if err != nil {
		return fmt.Errorf("repo patch user err: %w", err)
	}

	return nil
}
