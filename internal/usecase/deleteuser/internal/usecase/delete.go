package usecase

import (
	"context"
	"fmt"
)

func (u UseCase) Delete(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("repo delete user err: %w", err)
	}

	return nil
}
