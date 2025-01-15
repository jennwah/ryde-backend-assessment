package usecase

import (
	"context"
	"fmt"
)

func (u UseCase) Create(ctx context.Context, userID, friendID string) error {
	err := u.repo.Create(ctx, userID, friendID)
	if err != nil {
		return fmt.Errorf("repo create user friend pair err: %w", err)
	}

	return nil
}
