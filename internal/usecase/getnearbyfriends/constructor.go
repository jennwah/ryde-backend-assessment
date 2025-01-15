package getnearbyfriends

import (
	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/internal/repository"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/internal/usecase"
)

func Create(db *sqlx.DB) usecase.UseCase {
	return usecase.New(repository.New(db))
}
