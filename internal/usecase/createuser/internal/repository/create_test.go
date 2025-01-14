package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

func TestCreate(t *testing.T) {
	query := `INSERT INTO users (name, date_of_birth, address, description, location) VALUES (?, ?, ?, ?, public.ST_SetSRID(public.ST_MakePoint(?, ?), 4326)) RETURNING id`

	tests := []struct {
		name         string
		inputUser    model.User
		mockBehavior func(mock sqlmock.Sqlmock)
		wantID       string
		wantErr      bool
	}{
		{
			name: "Success create user with INSERT",
			inputUser: model.User{
				Name:        "John Doe",
				DateOfBirth: "1990-01-01",
				Address:     "123 Main St",
				Location: model.Point{
					Latitude:  1.234,
					Longitude: 2.345,
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(
					query,
				))
				mock.ExpectQuery(regexp.QuoteMeta(
					query,
				)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow("test-user-id"),
					)
			},
			wantID:  "test-user-id",
			wantErr: false,
		},
		{
			name: "Error preparing INSERT statement",
			inputUser: model.User{
				Name:        "Jane Doe",
				DateOfBirth: "1980-12-12",
				Address:     "456 Main St",
				Location: model.Point{
					Latitude:  10.001,
					Longitude: 20.002,
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(
					query,
				)).
					WillReturnError(fmt.Errorf("prepare statement error"))
			},
			wantID:  "",
			wantErr: true,
		},
		{
			name: "Error execute INSERT statement",
			inputUser: model.User{
				Name:        "Jake Doe",
				DateOfBirth: "1970-07-07",
				Address:     "999 Mystery Lane",
				Location: model.Point{
					Latitude:  123.456,
					Longitude: -45.678,
				},
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(
					query,
				))
				mock.ExpectQuery(regexp.QuoteMeta(
					query,
				)).
					WillReturnError(fmt.Errorf("some query error"))
			},
			wantID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			tt.mockBehavior(mock)
			repo := New(sqlxDB)

			gotID, err := repo.Create(context.Background(), tt.inputUser)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantID, gotID)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
