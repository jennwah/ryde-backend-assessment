package usecase

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

func (u UseCase) Create(ctx context.Context, user model.User) (string, error) {
	userID, err := u.repo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("repo create user err: %w", err)
	}

	return userID, nil
}
