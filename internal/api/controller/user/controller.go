package user

import (
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/config"
	createuseruc "github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser"
	deleteuseruc "github.com/jennwah/ryde-backend-engineer/internal/usecase/deleteuser"
	getuseruc "github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser"
	patchuseruc "github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser"
)

type UseCases struct {
	CreateUserUseCase createUser
	GetUserUseCase    getUser
	PatchUserUseCase  patchUser
	DeleteUserUseCase deleteUser
}

type Controller struct {
	uc     UseCases
	logger *slog.Logger
	cfg    config.Config
}

func New(db *sqlx.DB, logger *slog.Logger, cfg config.Config) Controller {
	RegisterUserValidators()

	return Controller{
		uc: UseCases{
			CreateUserUseCase: createuseruc.Create(db),
			GetUserUseCase:    getuseruc.Create(db),
			PatchUserUseCase:  patchuseruc.Create(db),
			DeleteUserUseCase: deleteuseruc.Create(db),
		},
		logger: logger,
		cfg:    cfg,
	}
}
