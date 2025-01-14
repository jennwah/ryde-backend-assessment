package usecase

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
)

func (u UseCase) Get(ctx context.Context, id string) (model.User, error) {
	user, err := u.repo.Get(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("repo get user err: %w", err)
	}

	return user, nil
}
