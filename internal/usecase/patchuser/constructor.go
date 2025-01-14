package patchuser

import (
	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/internal/repository"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/internal/usecase"
)

func Create(db *sqlx.DB) usecase.UseCase {
	return usecase.New(repository.New(db))
}
