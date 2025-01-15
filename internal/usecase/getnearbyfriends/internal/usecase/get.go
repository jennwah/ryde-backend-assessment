package usecase

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
)

func (u UseCase) GetNearby(ctx context.Context, userID string, radiusMeter float64) (model.Friends, error) {
	friends, err := u.repo.Get(ctx, userID, radiusMeter)
	if err != nil {
		return model.Friends{}, fmt.Errorf("repo get nearby friends err: %w", err)
	}

	return friends, nil
}
