package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		mockBehavior func(mock sqlmock.Sqlmock)
		expectedErr  error
	}{
		{
			name:   "success - user deleted",
			userID: "some-valid-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    DELETE FROM users
                    WHERE id = $1
                `)).
					WithArgs("some-valid-uuid").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "not found - sql.ErrNoRows",
			userID: "non-existent-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    DELETE FROM users
                    WHERE id = $1
                `)).
					WithArgs("non-existent-uuid").
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: postgresql.ErrNotFound,
		},
		{
			name:   "generic db error",
			userID: "random-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`
                    DELETE FROM users
                    WHERE id = $1
                `)).
					WithArgs("random-uuid").
					WillReturnError(errors.New("some db error"))
			},
			expectedErr: fmt.Errorf("delete user DB error: some db error"),
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

			err = repo.Delete(context.Background(), tt.userID)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr.Error())
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
