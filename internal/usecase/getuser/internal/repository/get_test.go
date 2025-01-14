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
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		mockBehavior func(sqlmock.Sqlmock)
		wantUser     model.User
		wantErr      error
	}{
		{
			name:   "success - user found",
			userID: "test-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
                    SELECT 
                        id,
                        name,
                        date_of_birth,
                        address,
                        description,
                        public.ST_Y(location::public.geometry) AS latitude,
                        public.ST_X(location::public.geometry) AS longitude
                    FROM users
                    WHERE id = $1
                `)).
					WithArgs("test-uuid").
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "date_of_birth", "address", "description", "latitude", "longitude"}).
							AddRow("test-uuid", "Test User", "2000-01-01", "123 Street", nil, 1.234, 2.345),
					)
			},
			wantUser: model.User{
				ID:          "test-uuid",
				Name:        "Test User",
				DateOfBirth: "2000-01-01",
				Address:     "123 Street",
				Description: nil,
				Location: model.Point{
					Latitude:  1.234,
					Longitude: 2.345,
				},
			},
			wantErr: nil,
		},
		{
			name:   "not found - no rows",
			userID: "missing-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
                    SELECT 
                        id,
                        name,
                        date_of_birth,
                        address,
                        description,
                        public.ST_Y(location::public.geometry) AS latitude,
                        public.ST_X(location::public.geometry) AS longitude
                    FROM users
                    WHERE id = $1
                `)).
					WithArgs("missing-uuid").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser: model.User{},
			wantErr:  postgresql.ErrNotFound,
		},
		{
			name:   "db error",
			userID: "random-uuid",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
                    SELECT 
                        id,
                        name,
                        date_of_birth,
                        address,
                        description,
                        public.ST_Y(location::public.geometry) AS latitude,
                        public.ST_X(location::public.geometry) AS longitude
                    FROM users
                    WHERE id = $1
                `)).
					WithArgs("random-uuid").
					WillReturnError(errors.New("some db error"))
			},
			wantUser: model.User{},
			wantErr:  fmt.Errorf("get user by id db err: some db error"),
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

			gotUser, gotErr := repo.Get(context.Background(), tt.userID)

			if tt.wantErr == nil {
				require.NoError(t, gotErr)
			} else {
				require.Error(t, gotErr)
				require.EqualError(t, gotErr, tt.wantErr.Error())
			}

			require.Equal(t, tt.wantUser, gotUser)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
