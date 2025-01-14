package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"strings"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

func (r *Repository) Create(ctx context.Context, user model.User) (string, error) {
	dbPayload := dbUser{
		Name:        user.Name,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		Description: user.Description,
		Latitude:    user.Location.Latitude,
		Longitude:   user.Location.Longitude,
	}

	// Create new user record
	stmt, err := r.db.PrepareNamedContext(
		ctx,
		`INSERT INTO users (name, date_of_birth, address, description, location)
		VALUES (:name, :date_of_birth, :address, :description, public.ST_SetSRID(public.ST_MakePoint(:longitude, :latitude), 4326)) 
		RETURNING id`,
	)
	if err != nil {
		return "", fmt.Errorf("failed prepare create user record, err: %w", err)
	}

	var userID string
	err = stmt.QueryRowContext(ctx, dbPayload).Scan(&userID)
	defer stmt.Close()

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) &&
			e.SQLState() == postgresql.ErrUniqueConstraintViolationCode &&
			strings.Contains(e.Message, "users_name_key") {
			return "", postgresql.ErrUserAlreadyExist
		}

		return "", fmt.Errorf("failed execute create user record, err: %w", err)
	}

	return userID, nil
}
