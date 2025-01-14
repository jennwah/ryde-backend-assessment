package repository

import (
	"context"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
)

func (r *Repository) Patch(ctx context.Context, payload model.PatchUser) error {
	dbPayload := dbUser{
		ID:          payload.ID,
		Name:        payload.Name,
		DateOfBirth: payload.DateOfBirth,
		Address:     payload.Address,
		Description: payload.Description,
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
	}

	query := `UPDATE users
		SET name = :name,
		    date_of_birth = :date_of_birth,
		    address = :address,
		    description = :description,
		    location = public.ST_SetSRID(public.ST_MakePoint(:longitude, :latitude), 4326)
		WHERE id = :id
	`

	res, err := r.db.NamedExecContext(ctx, query, dbPayload)
	if err != nil {
		return fmt.Errorf("update user DB error: %w", err)
	}
	row, _ := res.RowsAffected()
	if row == 0 {
		return postgresql.ErrNotFound
	}

	return nil
}
