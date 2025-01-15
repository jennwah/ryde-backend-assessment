package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"strings"
)

func (r *Repository) Create(ctx context.Context, userID, friendID string) error {
	pair := userFriendPair{}
	if userID < friendID {
		pair.UserID = userID
		pair.FriendID = friendID
	} else if friendID < userID {
		pair.UserID = friendID
		pair.FriendID = userID
	} else {
		return errors.New("user can't friend him/herself")
	}

	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO friends
		(user_id, friend_id)
		VALUES (:user_id, :friend_id)
		`,
		pair,
	)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) &&
			e.SQLState() == postgresql.ErrUniqueConstraintViolationCode &&
			strings.Contains(e.Message, "friends_pkey") {
			return postgresql.ErrFriendshipAlreadyExist
		}
		return fmt.Errorf("db insert user friend pair err: %w", err)
	}

	return nil
}
