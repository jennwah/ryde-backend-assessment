package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return postgresql.ErrNotFound
		}
		return fmt.Errorf("delete user DB error: %w", err)
	}

	return nil
}
