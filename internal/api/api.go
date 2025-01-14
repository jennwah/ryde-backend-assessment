package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	userCtrl "github.com/jennwah/ryde-backend-engineer/internal/api/controller/user"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

type API struct {
	userController userCtrl.Controller
}

func New(
	db *sqlx.DB,
	logger *slog.Logger,
	cfg config.Config,
) *API {
	return &API{
		userController: userCtrl.New(db, logger, cfg),
	}
}

func (a API) RegisterHandlers(r gin.IRouter) {
	userV1ApiGroup := r.Group("api/v1")

	userV1ApiGroup.GET("/hello", a.userController.Get)
}
