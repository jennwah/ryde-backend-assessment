package friend

import (
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/config"
	createfrienduc "github.com/jennwah/ryde-backend-engineer/internal/usecase/createfriend"
	getnearbyfriendsuc "github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends"
)

type UseCases struct {
	CreateFriendUseCase     createFriend
	GetNearbyFriendsUseCase getNearbyFriends
}

type Controller struct {
	uc     UseCases
	logger *slog.Logger
	cfg    config.Config
}

func New(db *sqlx.DB, logger *slog.Logger, cfg config.Config) Controller {
	return Controller{
		uc: UseCases{
			CreateFriendUseCase:     createfrienduc.Create(db),
			GetNearbyFriendsUseCase: getnearbyfriendsuc.Create(db),
		},
		logger: logger,
		cfg:    cfg,
	}
}
