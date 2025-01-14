package deleteuser

import (
	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/deleteuser/internal/repository"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/deleteuser/internal/usecase"
)

func Create(db *sqlx.DB) usecase.UseCase {
	return usecase.New(repository.New(db))
}
