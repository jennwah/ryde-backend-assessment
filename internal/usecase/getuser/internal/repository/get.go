package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
)

func (r *Repository) Get(ctx context.Context, id string) (model.User, error) {
	var dest dbUser

	query := `SELECT 
				id,
				name,
				date_of_birth,
				address,
				description,
				public.ST_Y(location::public.geometry) AS latitude,
				public.ST_X(location::public.geometry) AS longitude
			FROM users
			WHERE id = $1
	`

	if err := r.db.GetContext(ctx, &dest, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, postgresql.ErrNotFound
		}
		return model.User{}, fmt.Errorf("get user by id db err: %w", err)
	}

	return dest.toModel(), nil
}
