package repository

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
)

func (r *Repository) Get(ctx context.Context, userID string, radiusMeter float64) (model.Friends, error) {
	query := `
		WITH user_location AS (
			SELECT location
			FROM users.users
			WHERE id = $1
		)
		SELECT 
			u.id AS id,
			u.name,
			u.date_of_birth,
			u.address,
			u.description,
			public.ST_X(u.location::public.geometry) AS longitude,
			public.ST_Y(u.location::public.geometry) AS latitude
		FROM friends f
		JOIN users u ON 
			(f.user_id = $1 AND f.friend_id = u.id) OR
			(f.friend_id = $1 AND f.user_id = u.id)
		WHERE public.ST_DWithin(
			u.location,
			(SELECT location FROM user_location),
			$2 -- Radius in meters
		);
	`

	var friends dbFriends
	err := r.db.SelectContext(ctx, &friends, query, userID, radiusMeter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nearby friends: %w", err)
	}

	return friends.toModel(), nil
}
