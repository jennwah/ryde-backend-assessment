package user

import (
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

type UseCases struct {
	HelloUseCase helloWorld
}

type Controller struct {
	uc     UseCases
	logger *slog.Logger
	cfg    config.Config
}

func New(db *sqlx.DB, logger *slog.Logger, cfg config.Config) Controller {
	return Controller{
		uc: UseCases{
			HelloUseCase: nil,
		},
		logger: logger,
		cfg:    cfg,
	}
}
