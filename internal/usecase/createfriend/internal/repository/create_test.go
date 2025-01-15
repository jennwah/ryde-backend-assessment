package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		friendID     string
		mockBehavior func(mock sqlmock.Sqlmock)
		wantErr      error
	}{
		{
			name:     "success - userID < friendID",
			userID:   "aaa-000",
			friendID: "zzz-999",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO friends
                    (user_id, friend_id)
                    VALUES (?, ?)
                `)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name:     "success - friendID < userID",
			userID:   "zzz-999",
			friendID: "aaa-000",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO friends
                    (user_id, friend_id)
                    VALUES (?, ?)
                `)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name:         "fail - same user IDs",
			userID:       "same-id",
			friendID:     "same-id",
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      errors.New("user can't friend him/herself"),
		},
		{
			name:     "fail - unique constraint violation",
			userID:   "aaa-000",
			friendID: "zzz-999",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				pgErr := &pgconn.PgError{
					Code:    postgresql.ErrUniqueConstraintViolationCode,
					Message: "duplicate key value violates unique constraint \"friends_pkey\"",
				}
				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO friends
                    (user_id, friend_id)
                    VALUES (?, ?)
                `)).
					WillReturnError(pgErr)
			},
			wantErr: postgresql.ErrFriendshipAlreadyExist,
		},
		{
			name:     "fail - generic db error",
			userID:   "aaa-000",
			friendID: "zzz-999",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    INSERT INTO friends
                    (user_id, friend_id)
                    VALUES (?, ?)
                `)).
					WillReturnError(errors.New("some random db error"))
			},
			wantErr: fmt.Errorf("db insert user friend pair err: some random db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			repo := New(sqlxDB)
			tt.mockBehavior(mock)

			err = repo.Create(context.Background(), tt.userID, tt.friendID)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
