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
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
)

func TestPatch(t *testing.T) {
	query := `UPDATE users SET name = ?, date_of_birth = ?, address = ?, description = ?, location = public.ST_SetSRID(public.ST_MakePoint(?, ?), 4326) WHERE id = ?`

	tests := []struct {
		name         string
		patchPayload model.PatchUser
		mockBehavior func(sqlmock.Sqlmock)
		wantErr      error
	}{
		{
			name: "success - user patched",
			patchPayload: model.PatchUser{
				ID:          "test-uuid",
				Name:        "New Name",
				DateOfBirth: "1990-01-01",
				Address:     "123 Patch St.",
				Latitude:    1.234,
				Longitude:   2.345,
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "no rows found (not found)",
			patchPayload: model.PatchUser{
				ID:          "missing-uuid",
				Name:        "Name",
				DateOfBirth: "1970-01-01",
				Address:     "No Where",
				Latitude:    5.555,
				Longitude:   6.666,
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: postgresql.ErrNotFound,
		},
		{
			name: "generic db error",
			patchPayload: model.PatchUser{
				ID:          "some-uuid",
				Name:        "Broken Name",
				DateOfBirth: "1980-12-12",
				Address:     "Broken Address",
				Latitude:    10.111,
				Longitude:   20.222,
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("some db error"))
			},
			wantErr: fmt.Errorf("update user DB error: some db error"),
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

			err = repo.Patch(context.Background(), tt.patchPayload)

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
